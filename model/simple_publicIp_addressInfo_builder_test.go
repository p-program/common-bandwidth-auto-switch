package model

import (
	"testing"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/stretchr/testify/assert"
)

func TestBuild(t *testing.T) {
	new := prepareSimplePublicIpAddressInfoArray()
	t.Logf("new: %v", new)
	assert.Equal(t, len(new), 3)
}

func prepareSimplePublicIpAddressInfoArray() SimplePublicIpAddressInfos {
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
	entity1 := NewSimplePublicIpAddressInfoBuilder().
		AddPublicIpAddresse(&old[0]).
		AddValue(0).
		Build()
	entity2 := NewSimplePublicIpAddressInfoBuilder().
		AddPublicIpAddresse(&old[1]).
		AddValue(1).
		Build()
	entity3 := NewSimplePublicIpAddressInfoBuilder().
		AddPublicIpAddresse(&old[2]).
		AddValue(-1).
		Build()
	new := SimplePublicIpAddressInfos{*entity1, *entity2, *entity3}
	return new
}
