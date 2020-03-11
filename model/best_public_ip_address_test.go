package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDynamic(t *testing.T) {
	bandwidthInfos := prepareEipAvgBandwidthInfos()
	bestIPs, err := NewBestPublicIpAddress(40, bandwidthInfos)
	assert.Nil(t, err)
	bestIPs.dynamic()
}
func prepareEipAvgBandwidthInfos() []EipAvgBandwidthInfo {
	// bandwidthInfos := []EipAvgBandwidthInfo{
	// 	EipAvgBandwidthInfo{"1.1.1.1", "a", float64(20)},
	// 	EipAvgBandwidthInfo{"1.1.1.2", "", float64(10)},
	// 	EipAvgBandwidthInfo{"1.1.1.3", "", float64(15)},
	// 	EipAvgBandwidthInfo{"1.1.1.4", "", float64(5)},
	// 	EipAvgBandwidthInfo{"1.1.1.5", "", float64(2)},
	// 	EipAvgBandwidthInfo{"1.1.1.6", "", float64(21)},
	// }
	bandwidthInfos := []EipAvgBandwidthInfo{
		EipAvgBandwidthInfo{"1.1.1.1", "a", float64(21)},
		EipAvgBandwidthInfo{"1.1.1.2", "", float64(20)},
		EipAvgBandwidthInfo{"1.1.1.3", "", float64(31)},
	}
	return bandwidthInfos
}
