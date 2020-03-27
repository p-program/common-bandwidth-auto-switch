package manager

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
)

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

func TestScaleDown(t *testing.T) {
	eipAvgList := []model.EipAvgBandwidthInfo{
		{"1", "eip-1", float64(0.0869915403168777)},
		{"", "eip-", float64(0.006793844288793103)},
		{"", "eip-", float64(0.00014785240436422415)},
		{"", "eip-", float64(8.065813130345838)},
		{"", "eip-", float64(0.006415794635641164)},
		{"", "eip-", float64(40.98798738676926)},
		{"", "eip-", float64(0.001206628207502694)},
		{"", "eip-", float64(0.00011799253266433189)},
		{"", "eip-", float64(2.3567779804098197)},
		{"", "eip-", float64(0.22648515372440733)},
		{"", "eip-", float64(8.321491241455078)},
		{"", "eip-", float64(0.16179972681505927)},
		{"", "eip-", float64(0.39273702687230605)},
		{"", "eip-", float64(0.0001419330465382543)},
		{"", "eip-", float64(0.0002052044046336207)},
		{"", "eip-", float64(8.570368931211274)},
		{"", "eip-", float64(8.32183864198882)},
		{"", "eip-", float64(1.1245080355940194)},
		{"", "eip-", float64(7.19530829067888e-05)},
		{"", "eip-", float64(0.0001991535055226293)},
		{"", "eip-", float64(0.7311198793608567)},
		{"", "eip-", float64(0.1451458108836207)},
		{"", "eip-", float64(1.2695175697063577)},
		{"", "eip-", float64(0.07518465765591326)},
		{"", "eip-", float64(0.00015390330347521552)},
		{"", "eip-", float64(0.00010812693628771552)},
		{"", "eip-", float64(0.00025821554249730604)},
		{"", "eip-", float64(0.12906751961543642)},
		{"", "eip-", float64(0.12866618715483566)},
		{"", "eip-", float64(0.0817804007694639)},
		{"", "eip-", float64(2.265821128055967)},
		{"", "eip-", float64(0.03131892763335129)},
		{"", "eip-", float64(0.9026491888638201)},
		{"", "eip-", float64(0.8933019966914736)},
		{"", "eip-", float64(0.6154518127441406)},
		{"", "eip-", float64(0.009045568005792025)},
		{"", "eip-", float64(0.06862140523976293)},
		{"", "eip-", float64(0.039022774531923494)},
		{"", "eip-", float64(0.0019256328714304957)},
		{"", "eip-", float64(35.52917020074253)},
		// {"", "eip-", float64()},
		// {"", "eip-", float64()},
		// {"", "eip-", float64()},
		// {"", "eip-", float64()},
		// {"", "eip-", float64()},
		// {"", "eip-", float64()},
	}
	bestPublicIpAddress, err := model.NewBestPublicIpAddress(40, eipAvgList)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	bestIPs := bestPublicIpAddress.FindBestWithoutBrain()
	t.Log(bestIPs)
}

func prepareSDK(config *model.ProjectConfig) *aliyun.AliyunSDK {
	aliyunSDKConfig := config.AliyunConfig
	return aliyun.NewAliyunSDK(&aliyunSDKConfig)
}

func prepareConfing(t *testing.T) *model.ProjectConfig {
	config := model.NewProjectConfig()
	path := path.Join("../", "config.yaml")
	err := config.LoadYAML(path)
	assert.Nil(t, err)
	return config
}
