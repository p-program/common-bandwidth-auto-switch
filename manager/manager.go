package manager

import (
	"errors"
	"fmt"
	"math"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
	"github.com/zeusro/common-bandwidth-auto-switch/util"
)

const (
	ADD_EIP_TEMPLATE    = "添加EIP: %s ;EIPID: %s"
	REMOVE_EIP_TEMPLATE = "删除EIP: %s ;EIPID: %s"
)

// Manager 控制终端
type Manager struct {
	sdk                 *aliyun.AliyunSDK
	cbp                 *model.CommonBandwidthPackage
	dingtalkNotifyToken string
}

func NewManager(sdk *aliyun.AliyunSDK, cbp *model.CommonBandwidthPackage) *Manager {
	m := &Manager{
		sdk: sdk,
		cbp: cbp,
	}
	return m
}

// UseDingTalkNotification 使用钉钉消息推送
func (m *Manager) UseDingTalkNotification(token string) {
	m.dingtalkNotifyToken = token
}

func (m *Manager) Run() {
	var currentRateWG sync.WaitGroup
	currentRateWG.Add(2)
	//流入带宽检测数据平均值
	rxDataPoint := &model.Datapoint{}
	//流出带宽检测数据平均值
	txDataPoint := &model.Datapoint{}
	var err1, err2 error
	go func(w *sync.WaitGroup) {
		rxDataPoint, err1 = m.sdk.GetAvgRxRate(m.cbp)
		currentRateWG.Done()
	}(&currentRateWG)
	go func(w *sync.WaitGroup) {
		txDataPoint, err2 = m.sdk.GetAvgTxRate(m.cbp)
		currentRateWG.Done()
	}(&currentRateWG)
	currentRateWG.Wait()
	// 共享带宽信息
	cbpInfo := m.cbp
	log.Info().Msgf("当前共享带宽实例: %s ;平均流入带宽: %v ;平均流出带宽: %v ;",
		cbpInfo.ID,
		rxDataPoint.Value,
		txDataPoint.Value)
	errs := make([]error, 2)
	if err1 != nil {
		errs = append(errs, err1)
	}
	if err2 != nil {
		errs = append(errs, err2)
	}
	if len(errs) > 0 {
		log.Debug().Errs("goroutine get rate err", []error{err1, err2})
	}
	// 当前共享带宽最大带宽速率,单位是Mbps
	currentMaxBandwidthRate := math.Max(rxDataPoint.Value, txDataPoint.Value)
	if currentMaxBandwidthRate > float64(cbpInfo.MaxBandwidth) {
		//带宽高峰，需要缩容
		m.ScaleDown(currentMaxBandwidthRate)
		return
	}
	if float64(cbpInfo.MinBandwidth)-currentMaxBandwidthRate > 5 {
		//带宽低谷，需要扩容
		m.ScaleUp(currentMaxBandwidthRate)
		return
	}
	// 5 Mbps 以内就不优化了，没啥区别
	//无需扩容,也无需缩容
	log.Info().Msg("无需扩容,也无需缩容")
}

// ScaleUp 扩容:将低带宽EIP加入共享带宽
func (m *Manager) ScaleUp(currentBandwidthRate float64) (err error) {
	cbpInfo := m.cbp
	cbwpID := cbpInfo.ID
	//获取当前 region 下未绑定共享带宽的IP列表
	currentUnbindEIPs, err := m.sdk.GetCurrentEipAddressesExceptCBWP(cbwpID)
	if err != nil {
		return err
	}
	var eipWaitLock sync.WaitGroup
	checkFrequency := cbpInfo.CheckFrequency
	eipWaitLock.Add(len(currentUnbindEIPs))
	var eipAvgList []model.EipAvgBandwidthInfo
	for _, eipInfo := range currentUnbindEIPs {
		go func(eip *vpc.EipAddress, wg *sync.WaitGroup) {
			defer wg.Done()
			avgBandwidth, err := m.sdk.DescribeEipAvgMonitorData(eip.AllocationId, checkFrequency)
			//FIXME: 局部失败要怎么处理
			if err != nil {
				log.Err(err)
				return
			}
			eipAvgList = append(eipAvgList, model.EipAvgBandwidthInfo{
				IpAddress:    eip.IpAddress,
				AllocationId: eip.AllocationId,
				Value:        avgBandwidth,
			})
		}(&eipInfo, &eipWaitLock)
	}
	//根据剩余带宽动态规划
	bandwidthLimit := m.cbp.MinBandwidth - int(currentBandwidthRate)
	bestPublicIpAddress, err := model.NewBestPublicIpAddress(bandwidthLimit, eipAvgList)
	if err != nil {
		return err
	}
	// 动态优化求最优IP
	bestEIPs := bestPublicIpAddress.FindBestWithoutBrain()
	if len(bestEIPs) < 1 {
		return nil
	}
	m.ding(bestEIPs, ADD_EIP_TEMPLATE)
	// TODO: 周密测试后再取消注释
	// for _, eipInfo := range bestEIPs {
	// 	m.sdk.AddCommonBandwidthPackageIp(cbpInfo.ID, eipInfo.AllocationId)
	// }
	return nil
}

func (m *Manager) ding(ips []model.EipAvgBandwidthInfo, notifyTemplate string) {
	cbwpID := m.cbp.ID
	if token := m.dingtalkNotifyToken; len(token) > 0 {
		ding := util.NewDingTalk(token)
		markdownBuilder := util.NewMarkdownBuilder()
		for _, eipInfo := range ips {
			content := fmt.Sprintf(ADD_EIP_TEMPLATE, eipInfo.IpAddress, eipInfo.AllocationId)
			markdownBuilder.AddText(content)
		}
		title := fmt.Sprintf("共享带宽动态优化(%s)", cbwpID)
		ding.DingMarkdown(title, markdownBuilder.BuilderText())
	}
}

//ScaleDown 缩容:将高带宽EIP移除共享带宽
func (m *Manager) ScaleDown(currentBandwidthRate float64) (err error) {
	cbpInfo := m.cbp
	// 获取当前共享带宽内EIP列表
	eipList, err := m.sdk.DescribeCommonBandwidthPackages(cbpInfo.ID)
	if err != nil {
		return err
	}
	if len(eipList) == 0 {
		return errors.New("len(eipList) == 0")
	}
	//获取所有EIP监控
	var eipWaitLock sync.WaitGroup
	eipWaitLock.Add(len(eipList))
	checkFrequency := cbpInfo.CheckFrequency
	var eipAvgList []model.EipAvgBandwidthInfo
	for _, eipInfo := range eipList {
		go func(eip *vpc.PublicIpAddresse, wg *sync.WaitGroup) {
			defer wg.Done()
			avgBandwidth, err := m.sdk.DescribeEipAvgMonitorData(eip.AllocationId, checkFrequency)
			//FIXME: 局部失败要怎么处理
			if err != nil {
				log.Err(err)
				return
			}
			eipAvgList = append(eipAvgList, model.EipAvgBandwidthInfo{
				IpAddress:    eip.IpAddress,
				AllocationId: eip.AllocationId,
				Value:        avgBandwidth,
			})
		}(&eipInfo, &eipWaitLock)
	}
	bestPublicIpAddress, err := model.NewBestPublicIpAddress(m.cbp.MinBandwidth, eipAvgList)
	if err != nil {
		return err
	}
	//进行动态优化
	bestIPs := bestPublicIpAddress.FindBestWithoutBrain()
	if len(bestIPs) < 1 {
		return nil
	}
	currentEIPsInCBWP, err := m.sdk.GetCurrentEipAddressesInCBWP(cbpInfo.ID)
	if err != nil {
		return err
	}
	var lowestEIPs []model.EipAvgBandwidthInfo
	var lowestEIPsAddress []string
	//求差集. currentEIPsInCBWP - bestIPs
	for i1 := 0; i1 < len(bestIPs); i1++ {
		bestIP := bestIPs[i1]
		isCross := false
		var ip vpc.EipAddress
		for i2 := 0; i2 < len(currentEIPsInCBWP); i2++ {
			ip = currentEIPsInCBWP[i2]
			if bestIP.AllocationId == ip.AllocationId {
				isCross = true
				break
			}
		}
		//没交集,可加
		if !isCross {
			entity := model.EipAvgBandwidthInfo{
				IpAddress:    ip.IpAddress,
				AllocationId: ip.AllocationId,
			}
			lowestEIPs = append(lowestEIPs, entity)
			lowestEIPsAddress = append(lowestEIPsAddress, ip.AllocationId)
		}
	}
	m.ding(lowestEIPs, REMOVE_EIP_TEMPLATE)

	// TODO: 周密测试后再取消注释
	// m.sdk.RemoveCommonBandwidthPackageIps(cbpInfo.ID, lowestEIPsAddress)
	return nil
}
