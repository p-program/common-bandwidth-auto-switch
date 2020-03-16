package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

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
	configPath := path.Join(homeDir, defaultConfig)
	err = config.LoadYAML(configPath)
	if err != nil {
		log.Warn().Msg(err.Error())
		configPath = path.Join(homeDir, exampleConfig)
		err := config.LoadYAML(configPath)
		if err != nil {
			panic(err)
		}
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
	fmt.Println("readFile")
	readFile()

}

func readFile() {
	files, err := ioutil.ReadDir("/app2")
	if err != nil {
		log.Err(err)
		fmt.Print("/app2 has no file")
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
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
