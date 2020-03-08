package aliyun

import (
	"errors"
	"fmt"
	"time"

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

// DescribeEipMonitorData 调用DescribeEipMonitorData接口查看EIP的监控信息。
// https://help.aliyun.com/document_detail/36060.html
func (sdk *AliyunSDK) DescribeEipMonitorData(allocationId string, checkFrequency string) ([]vpc.EipMonitorData, error) {
	client := sdk.GetVPCClient()
	request := vpc.CreateDescribeEipMonitorDataRequest()
	request.Scheme = "https"
	// 60 秒一个周期
	request.Period = "60"
	now := time.Now()
	frequency, err := time.ParseDuration(checkFrequency)
	log.Info().Msgf("duration: %s", frequency.String())
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
