package aliyun

import (
	"errors"
	"fmt"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/rs/zerolog/log"
)

// GetPublicIpAddresseFromCommonBandwidthPackages 从第一个共享带宽里面提取EIP信息
func (sdk *AliyunSDK) GetPublicIpAddresseFromCommonBandwidthPackages() (eipList []vpc.PublicIpAddresse, err error) {
	client, err := vpc.NewClientWithAccessKey(sdk.config.Region, sdk.config.AccessKeyId, sdk.config.AccessSecret)
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
