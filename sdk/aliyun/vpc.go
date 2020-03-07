package aliyun

import (
	"errors"
	"fmt"

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

// GetPublicIpAddresseFromCommonBandwidthPackages 从第一个共享带宽里面提取EIP信息
func (sdk *AliyunSDK) GetPublicIpAddresseFromCommonBandwidthPackages() (eipList []vpc.PublicIpAddresse, err error) {
	client := sdk.GetVPCClient()
	request := vpc.CreateDescribeCommonBandwidthPackagesRequest()
	request.Scheme = "https"
	request.PageSize = "50"
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
	// request.RegionId = sdk.config.Region
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
