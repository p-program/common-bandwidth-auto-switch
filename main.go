package main

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"strings"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/manager"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
)

const (
	defaultConfig = "config.yaml"
	exampleConfig = "config-example.yaml"
	myName        = `
✄╔════╗
✄╚══╗═║
✄──╔╝╔╝╔══╗╔╗╔╗╔══╗╔═╗╔══╗
✄─╔╝╔╝─║║═╣║║║║║══╣║╔╝║╔╗║
✄╔╝═╚═╗║║═╣║╚╝║╠══║║║─║╚╝║
✄╚════╝╚══╝╚══╝╚══╝╚╝─╚══╝
`
	LINE = "----------------------------------------"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

func main() {
	fmt.Println(LINE)
	fmt.Print("Power by")
	fmt.Println(myName)
	fmt.Println(LINE)
	setMaxProcs()
	config := model.NewProjectConfig()
	homeDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configDir := path.Join(homeDir, "config")
	configPath := path.Join(configDir, defaultConfig)
	err = config.LoadYAML(configPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		configPath = path.Join(configDir, exampleConfig)
		err := config.LoadYAML(configPath)
		if err != nil {
			panic(err)
		}
	}
	if strings.EqualFold("debug", config.LogLevel) {
		log.Info().Msgf("config.LogLevel:%s", config.LogLevel)
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	aliyunSDKConfig := config.AliyunConfig
	useDingTalkNotification := len(config.DingTalkConfig.NotificationToken) > 0
	sdk := aliyun.NewAliyunSDK(&aliyunSDKConfig)
	for _, cbp := range config.CommonBandwidthPackages {
		manager := manager.NewManager(sdk, &cbp)
		if useDingTalkNotification {
			manager.UseDingTalkNotification(config.DingTalkConfig.NotificationToken)
		}
		manager.Run()
	}
}

func setMaxProcs() {
	// Allow as many threads as we have cores unless the user specified a value.
	numProcs := runtime.NumCPU()
	runtime.GOMAXPROCS(numProcs)
	// Check if the setting was successful.
	actualNumProcs := runtime.GOMAXPROCS(0)
	if actualNumProcs != numProcs {
		log.Info().Msgf("Specified max procs of %d but using %d", numProcs, actualNumProcs)
	}
}
