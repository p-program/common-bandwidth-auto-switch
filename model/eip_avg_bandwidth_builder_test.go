package model

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	new := prepareEipAvgBandwidthInfoArray()
	t.Logf("new: %v", new)
	assert.Equal(t, len(new), 3)
}

func prepareEipAvgBandwidthInfoArray() EipAvgBandwidthInfos {
	old := []vpc.PublicIpAddresse{
		vpc.PublicIpAddresse{
			IpAddress:    "1.1.1.1",
			AllocationId: "a",
		},
		vpc.PublicIpAddresse{
			IpAddress:    "1.1.1.3",
			AllocationId: "c",
		},
		vpc.PublicIpAddresse{
			IpAddress:    "1.1.1.2",
			AllocationId: "b",
		},
	}
	entity1 := NewEipAvgBandwidthInfoBuilder().
		AddPublicIpAddresse(&old[0]).
		AddValue(0).
		Build()
	entity2 := NewEipAvgBandwidthInfoBuilder().
		AddPublicIpAddresse(&old[1]).
		AddValue(1).
		Build()
	entity3 := NewEipAvgBandwidthInfoBuilder().
		AddPublicIpAddresse(&old[2]).
		AddValue(-1).
		Build()
	new := EipAvgBandwidthInfos{*entity1, *entity2, *entity3}
	return new
}
