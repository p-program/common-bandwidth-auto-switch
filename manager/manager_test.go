package manager

import (
	"os"
	"path"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func TestRun(t *testing.T) {
	cfg := prepareConfing(t)
	useDingTalkNotification := len(cfg.DingTalkConfig.NotificationToken) > 0
	sdk := prepareSDK(cfg)
	for _, cbp := range cfg.CommonBandwidthPackages {
		manager := NewManager(sdk, &cbp)
		if useDingTalkNotification {
			manager.UseDingTalkNotification(cfg.DingTalkConfig.NotificationToken)
		}
		manager.Run()
	}
}

func prepareSDK(config *model.ProjectConfig) *aliyun.AliyunSDK {
	aliyunSDKConfig := config.AliyunConfig
	return aliyun.NewAliyunSDK(&aliyunSDKConfig)
}

func prepareConfing(t *testing.T) *model.ProjectConfig {
	config := model.NewProjectConfig()
	path := path.Join("../", "config", "config.yaml")
	err := config.LoadYAML(path)
	assert.Nil(t, err)
	return config
}
