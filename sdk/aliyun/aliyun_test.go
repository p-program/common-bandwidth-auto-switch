package aliyun

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
)

func prepareSDK(config *model.ProjectConfig) *AliyunSDK {
	aliyunSDKConfig := config.AliyunConfig
	return NewAliyunSDK(&aliyunSDKConfig)
}

func prepareConfig(t *testing.T) *model.ProjectConfig {
	config := model.NewProjectConfig()
	path := path.Join("../", "../", "config.yaml")
	err := config.LoadYAML(path)
	assert.Nil(t, err)
	return config
}
