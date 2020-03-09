package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindBestPublicIpAddress(t *testing.T) {
	bandwidthInfos := prepareEipAvgBandwidthInfos()
	bestIPs, err := NewBestPublicIpAddress(40, bandwidthInfos)
	assert.Nil(t, err)
	bestIPs.FindBestPublicIpAddress()
}
func prepareEipAvgBandwidthInfos() []EipAvgBandwidthInfo {
	bandwidthInfos := []EipAvgBandwidthInfo{
		EipAvgBandwidthInfo{"1.1.1.1", "a", float64(20)},
		EipAvgBandwidthInfo{"1.1.1.2", "", float64(10)},
		EipAvgBandwidthInfo{"1.1.1.3", "", float64(15)},
		EipAvgBandwidthInfo{"1.1.1.4", "", float64(5)},
		EipAvgBandwidthInfo{"1.1.1.5", "", float64(2)},
		EipAvgBandwidthInfo{"1.1.1.6", "", float64(21)},
	}
	return bandwidthInfos
}
