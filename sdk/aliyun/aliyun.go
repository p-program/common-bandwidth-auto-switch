package aliyun

import (
	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
)

type AliyunSDK struct {
	config *model.AliyunConfig
	client *vpc.Client
}

func NewAliyunSDK(config *model.AliyunConfig) *AliyunSDK {

	return &AliyunSDK{
		config: config,
	}
}

func (sdk AliyunSDK) GetClient() *vpc.Client {
	if sdk.client != nil {
		return sdk.client
	}
	client, err := vpc.NewClientWithAccessKey(sdk.config.Region, sdk.config.AccessKeyId, sdk.config.AccessSecret)
	if err != nil {
		log.Err(err)
		return nil
	}
	sdk.client = client
	return client
}
