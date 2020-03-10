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
	TEST_EIP_ID = "eip-xxx"
)

func TestDescribeEipMonitorData(t *testing.T) {
	config := prepareConfig(t)
	sdk := prepareSDK(config)
	cbs := config.CommonBandwidthPackages[0]
	list, err := sdk.DescribeEipMonitorData(TEST_EIP_ID, cbs.CheckFrequency)
	assert.Nil(t, err)
	for _, v := range list {
		t.Logf("EipBandwidth: %v ;EipFlow: %v ;流入带宽EipRX: %v ;流出带宽EipTX: %v ;", v.EipBandwidth, v.EipFlow, v.EipRX, v.EipTX)
	}
	// t.Logf("list:%v", list)
}

func TestDescribeEipAvgMonitorData(t *testing.T) {
	config := prepareConfig(t)
	sdk := prepareSDK(config)
	cbs := config.CommonBandwidthPackages[0]
	avgBandwidth, err := sdk.DescribeEipAvgMonitorData(TEST_EIP_ID, cbs.CheckFrequency)
	assert.Nil(t, err)
	t.Logf("EIP ID: %s ; avgBandwidth: %v Mbps", TEST_EIP_ID, avgBandwidth)
}
