package aliyun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAvgRxRate(t *testing.T) {
	config := prepareConfig(t)
	aliyunSDKConfig := config.AliyunConfig
	sdk := NewAliyunSDK(&aliyunSDKConfig)
	dataPoint, err := sdk.GetAvgRxRate(config.Frequency)
	assert.Nil(t, err)
	t.Logf("%v", dataPoint)
}

func TestGetAvgTxRate(t *testing.T) {
	config := prepareConfig(t)
	aliyunSDKConfig := config.AliyunConfig
	sdk := NewAliyunSDK(&aliyunSDKConfig)
	dataPoint, err := sdk.GetAvgTxRate(config.Frequency)
	assert.Nil(t, err)
	t.Logf("%v", dataPoint)
}
