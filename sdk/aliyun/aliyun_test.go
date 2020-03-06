package aliyun

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
)

func prepareSDK(t *testing.T) *AliyunSDK {
	config := model.NewProjectConfig()
	path := path.Join("../", "../", "config.yaml")
	err := config.LoadYAML(path)
	assert.Nil(t, err)
	aliyunSDKConfig := config.AliyunConfig
	return NewAliyunSDK(&aliyunSDKConfig)
}
