package aliyun

import (
	"github.com/zeusro/common-bandwidth-auto-switch/model"
)

type AliyunSDK struct {
	config *model.AliyunConfig
}

func NewAliyunSDK(config *model.AliyunConfig) *AliyunSDK {
	return &AliyunSDK{
		config: config,
	}
}
