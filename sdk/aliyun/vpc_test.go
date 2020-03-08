package aliyun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIpAddresseFromCommonBandwidthPackages(t *testing.T) {
	config := prepareConfig(t)
	sdk := prepareSDK(config)
	cbs := config.CommonBandwidthPackages[0]
	eipList, err := sdk.DescribeCommonBandwidthPackages(cbs.ID)
	assert.Nil(t, err)
	assert.Greater(t, len(eipList), 0)
}

const (
	TEST_EIP_ID = "eip-xxxx"
)

func TestDescribeEipMonitorData(t *testing.T) {
	config := prepareConfig(t)
	sdk := prepareSDK(config)
	cbs := config.CommonBandwidthPackages[0]
	list, err := sdk.DescribeEipMonitorData(TEST_EIP_ID, cbs.CheckFrequency)
	assert.Nil(t, err)
	t.Logf("list:%v", list)
}
