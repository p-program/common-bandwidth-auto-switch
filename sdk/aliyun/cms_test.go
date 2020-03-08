package aliyun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAvgRxRate(t *testing.T) {
	config := prepareConfig(t)
	aliyunSDKConfig := config.AliyunConfig
	sdk := NewAliyunSDK(&aliyunSDKConfig)
	cbpInfo := &config.CommonBandwidthPackages[0]
	dataPoint, err := sdk.GetAvgTxRate(cbpInfo)
	assert.Nil(t, err)
	t.Logf("%v", dataPoint)
}

func TestGetAvgTxRate(t *testing.T) {
	config := prepareConfig(t)
	aliyunSDKConfig := config.AliyunConfig
	sdk := NewAliyunSDK(&aliyunSDKConfig)
	cbpInfo := &config.CommonBandwidthPackages[0]
	dataPoint, err := sdk.GetAvgTxRate(cbpInfo)
	assert.Nil(t, err)
	t.Logf("%v", dataPoint)
}
