package aliyun

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetPublicIpAddresseFromCommonBandwidthPackages(t *testing.T) {
	sdk := prepareSDK(t)
	eipList, err := sdk.GetPublicIpAddresseFromCommonBandwidthPackages()
	assert.Nil(t, err)
	assert.Greater(t, len(eipList), 0)
}
