package aliyun

import (
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/rs/zerolog/log"
)

// GetVPCClient 共享带宽相关接口
// https://help.aliyun.com/document_detail/55928.html?spm=a2c4g.11186623.6.572.545b7579DjZDj2
func (sdk *AliyunSDK) GetVPCClient() *vpc.Client {
	client, err := vpc.NewClientWithAccessKey(sdk.config.Region, sdk.config.AccessKeyId, sdk.config.AccessSecret)
	if err != nil {
		log.Err(err)
		return nil
	}
	return client
}

// DescribeCommonBandwidthPackages 提取共享带宽里面EIP信息
func (sdk *AliyunSDK) DescribeCommonBandwidthPackages(bandwidthPackageId string) (eipList []vpc.PublicIpAddresse, err error) {
	client := sdk.GetVPCClient()
	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	request.Scheme = "https"
	request.PageSize = "50"
	request.BandwidthPackageId = bandwidthPackageId
	response, err := client.DescribeCommonBandwidthPackages(request)
	if err != nil {
		fmt.Print(err.Error())
	}
	packages := response.CommonBandwidthPackages.CommonBandwidthPackage
	if len(packages) < 1 {
		err = errors.New("No response.CommonBandwidthPackages.CommonBandwidthPackage,此region下没有共享带宽")
		return
	}
	commonBandwidthPackage := packages[0]
	log.Info().Msgf("Name: %s; BandwidthPackageId: %s",
		commonBandwidthPackage.Name,
		commonBandwidthPackage.BandwidthPackageId)
	eipList = commonBandwidthPackage.PublicIpAddresses.PublicIpAddresse
	return eipList, nil
}

// AddCommonBandwidthPackageIp 调用AddCommonBandwidthPackageIp接口添加EIP到共享带宽中。
func (sdk *AliyunSDK) AddCommonBandwidthPackageIp(bandwidthPackageId string, ipInstanceId string) bool {
	client := sdk.GetVPCClient()
	request := vpc.CreateAddCommonBandwidthPackageIpRequest()
	request.Scheme = "https"
	request.BandwidthPackageId = bandwidthPackageId
	request.IpInstanceId = ipInstanceId
	response, err := client.AddCommonBandwidthPackageIp(request)
	if err != nil {
		log.Err(err)
		return false
	}
	return response.IsSuccess()
}

// RemoveCommonBandwidthPackageIp 调用RemoveCommonBandwidthPackageIp接口移除共享带宽实例中的EIP。
// https://help.aliyun.com/document_detail/55995.html
func (sdk *AliyunSDK) RemoveCommonBandwidthPackageIp(bandwidthPackageId string, ipInstanceId string) bool {
	client := sdk.GetVPCClient()
	request := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
	request.Scheme = "https"
	request.RegionId = sdk.config.Region
	request.BandwidthPackageId = bandwidthPackageId
	request.IpInstanceId = ipInstanceId
	response, err := client.RemoveCommonBandwidthPackageIp(request)
	if err != nil {
		log.Err(err)
		return false
	}
	return response.IsSuccess()
}

// RemoveCommonBandwidthPackageIps 复用 VPC client ，并行删除EIP
func (sdk *AliyunSDK) RemoveCommonBandwidthPackageIps(bandwidthPackageId string, ipInstanceIds []string) {
	client := sdk.GetVPCClient()
	request := vpc.CreateRemoveCommonBandwidthPackageIpRequest()
	request.Scheme = "https"
	request.RegionId = sdk.config.Region
	request.BandwidthPackageId = bandwidthPackageId
	//  并发删除会报错
	var wg sync.WaitGroup
	wg.Add(len(ipInstanceIds))
	for _, ipInstanceID := range ipInstanceIds {
		//这里传复制，防止出错
		go func(c vpc.Client, r vpc.RemoveCommonBandwidthPackageIpRequest, eipID string, w *sync.WaitGroup) {
			//无论失败与否都解除占用
			defer w.Done()
			r.IpInstanceId = eipID
			_, err := c.RemoveCommonBandwidthPackageIp(&r)
			if err != nil {
				log.Err(err)
			}
		}(*client, *request, ipInstanceID, &wg)
	}
	wg.Wait()
}

// DescribeEipAvgMonitorData 获取 EIP 监控信息流入和流出的带宽总和平均值
// avgBandwidth 单位是 Mbps
func (sdk *AliyunSDK) DescribeEipAvgMonitorData(allocationId string, checkFrequency string) (avgBandwidth float64, err error) {
	datas, err := sdk.DescribeEipMonitorData(allocationId, checkFrequency)
	if err != nil {
		log.Err(err)
		return 0, err
	}
	var sum float64 = 0
	for _, data := range datas {
		// EipFlow = 流入和流出的带宽总和
		// EipBandwidth 带宽值，该值等于EipFlow/60，单位为B/S
		sum += float64(data.EipBandwidth)
	}
	// 1 Mbps = 131072 B/S
	avgBandwidth = sum / float64(len(datas)*131072)
	return avgBandwidth, nil
}

// DescribeEipMonitorData 调用DescribeEipMonitorData接口查看EIP的监控信息。
// https://help.aliyun.com/document_detail/36060.html
// allocationId EIP的实例ID。
func (sdk *AliyunSDK) DescribeEipMonitorData(allocationId string, checkFrequency string) ([]vpc.EipMonitorData, error) {
	client := sdk.GetVPCClient()
	request := vpc.CreateDescribeEipMonitorDataRequest()
	request.Scheme = "https"
	// 60 秒一个周期
	request.Period = "60"
	now := time.Now()
	frequency, err := time.ParseDuration(checkFrequency)
	// log.Info().Msgf("duration: %s", frequency.String())
	if err != nil {
		return nil, err
	}
	apiRequiredFormat := "2006-01-02T15:04:05Z"
	request.StartTime = now.Add(-frequency).UTC().Format(apiRequiredFormat)
	request.EndTime = now.UTC().Format(apiRequiredFormat)
	request.AllocationId = allocationId
	response, err := client.DescribeEipMonitorData(request)
	if err != nil {
		return nil, err
	}
	if !response.IsSuccess() {
		err = errors.New(response.BaseResponse.String())
	}
	return response.EipMonitorDatas.EipMonitorData, nil
}

// DescribeInuseEipAddresses 获取当前 region 所有在用 EIP
// https://help.aliyun.com/document_detail/36018.html
func (sdk *AliyunSDK) DescribeInuseEipAddresses() ([]vpc.EipAddress, error) {
	client := sdk.GetVPCClient()
	request := vpc.CreateDescribeEipAddressesRequest()
	request.RegionId = sdk.config.Region
	request.Scheme = "https"
	request.Status = "InUse"
	request.PageSize = requests.NewInteger(100)
	response, err := client.DescribeEipAddresses(request)
	return response.EipAddresses.EipAddress, err
}

// GetCurrentEipAddressesExceptCBWP 获取未绑定共享带宽的EIP列表
// cbwpID 共享带宽ID
func (sdk *AliyunSDK) GetCurrentEipAddressesExceptCBWP(cbwpID string) (finalList []vpc.EipAddress, err error) {
	list, err := sdk.DescribeInuseEipAddresses()
	if err != nil {
		return
	}
	for _, item := range list {
		if !strings.EqualFold(cbwpID, item.BandwidthPackageId) {
			finalList = append(finalList, item)
		}
	}
	return finalList, nil
}

// GetCurrentEipAddressesInCBWP 获取绑定当前共享带宽的IP
func (sdk *AliyunSDK) GetCurrentEipAddressesInCBWP(cbwpID string) (finalList []vpc.EipAddress, err error) {
	list, err := sdk.DescribeInuseEipAddresses()
	if err != nil {
		return
	}
	for _, item := range list {
		if strings.EqualFold(cbwpID, item.BandwidthPackageId) {
			finalList = append(finalList, item)
		}
	}
	return finalList, nil
}
