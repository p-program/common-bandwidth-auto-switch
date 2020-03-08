package main

import (
	"fmt"
	"path"
	"runtime"

	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/manager"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
)

func main() {
	fmt.Println("start")
	fmt.Println("Power by Zeusro")
	setMaxProcs()
	config := model.NewProjectConfig()
	path := path.Join("config.yaml")
	err := config.LoadYAML(path)
	if err != nil {
		path = "config-example.yaml"
		err := config.LoadYAML(path)
		if err != nil {
			panic(err)
		}
	}
	aliyunSDKConfig := config.AliyunConfig
	sdk := aliyun.NewAliyunSDK(&aliyunSDKConfig)
	for _, cbp := range config.CommonBandwidthPackages {
		manager := manager.NewManager(sdk, &cbp)
		manager.Run()
	}

}

func prepareSDK() *aliyun.AliyunSDK {
	config := model.NewProjectConfig()
	path := path.Join("config.yaml")
	err := config.LoadYAML(path)
	if err != nil {
		path = "config-example.yaml"
		err := config.LoadYAML(path)
		if err != nil {
			panic(err)
		}
	}
	aliyunSDKConfig := config.AliyunConfig
	sdk := aliyun.NewAliyunSDK(&aliyunSDKConfig)
	return sdk
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
