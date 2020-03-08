package model

import "github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

type SimplePublicIpAddressInfoBuilder struct {
	publicIpAddressInfo *SimplePublicIpAddressInfo
	// publicIpAddressInfos []SimplePublicIpAddressInfo
}

func NewSimplePublicIpAddressInfoBuilder() *SimplePublicIpAddressInfoBuilder {
	return &SimplePublicIpAddressInfoBuilder{}
}

func (b *SimplePublicIpAddressInfoBuilder) Build() *SimplePublicIpAddressInfo {
	return b.publicIpAddressInfo
}

// func (b *SimplePublicIpAddressInfoBuilder) BuildList() []SimplePublicIpAddressInfo {
// 	return b.publicIpAddressInfos
// }

func (b *SimplePublicIpAddressInfoBuilder) AddPublicIpAddresse(old *vpc.PublicIpAddresse) *SimplePublicIpAddressInfoBuilder {
	if b.publicIpAddressInfo == nil {
		b.publicIpAddressInfo = &SimplePublicIpAddressInfo{
			IpAddress:    old.IpAddress,
			AllocationId: old.AllocationId,
		}
	} else {
		b.publicIpAddressInfo.IpAddress = old.IpAddress
		b.publicIpAddressInfo.AllocationId = old.AllocationId
	}
	return b
}

func (b *SimplePublicIpAddressInfoBuilder) AddValue(value float64) *SimplePublicIpAddressInfoBuilder {
	if b.publicIpAddressInfo == nil {
		b.publicIpAddressInfo = &SimplePublicIpAddressInfo{
			Value: value,
		}
	} else {
		b.publicIpAddressInfo.Value = value
	}
	return b
}
