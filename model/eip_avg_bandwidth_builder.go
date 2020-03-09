package model

import "github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"

type EipAvgBandwidthInfoBuilder struct {
	publicIpAddressInfo *EipAvgBandwidthInfo
}

func NewEipAvgBandwidthInfoBuilder() *EipAvgBandwidthInfoBuilder {
	return &EipAvgBandwidthInfoBuilder{}
}

func (b *EipAvgBandwidthInfoBuilder) Build() *EipAvgBandwidthInfo {
	return b.publicIpAddressInfo
}

func (b *EipAvgBandwidthInfoBuilder) AddPublicIpAddresse(old *vpc.PublicIpAddresse) *EipAvgBandwidthInfoBuilder {
	if b.publicIpAddressInfo == nil {
		b.publicIpAddressInfo = &EipAvgBandwidthInfo{
			IpAddress:    old.IpAddress,
			AllocationId: old.AllocationId,
		}
	} else {
		b.publicIpAddressInfo.IpAddress = old.IpAddress
		b.publicIpAddressInfo.AllocationId = old.AllocationId
	}
	return b
}

func (b *EipAvgBandwidthInfoBuilder) AddValue(value float64) *EipAvgBandwidthInfoBuilder {
	if b.publicIpAddressInfo == nil {
		b.publicIpAddressInfo = &EipAvgBandwidthInfo{
			Value: value,
		}
	} else {
		b.publicIpAddressInfo.Value = value
	}
	return b
}
