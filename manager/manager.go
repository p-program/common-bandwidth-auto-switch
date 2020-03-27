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
)

const (
	ADD_EIP_TEMPLATE    = "添加EIP: %s ;\n\nIP: %s "
	REMOVE_EIP_TEMPLATE = "删除EIP: %s ;\n\nIP: %s "
)

// Manager 控制终端
type Manager struct {
	sdk                 *aliyun.AliyunSDK
	cbp                 *model.CommonBandwidthPackage
	dingtalkNotifyToken string
	//dryRun true 表示以测试模式运行，最后不会真的添加/删除EIP
	dryRun bool
}

func NewManager(sdk *aliyun.AliyunSDK, cbp *model.CommonBandwidthPackage) *Manager {
	m := &Manager{
		sdk: sdk,
		cbp: cbp,
	}
	return m
}

// DryRun true 表示以测试模式运行，最后不会真的添加/删除EIP
func (m *Manager) DryRun(dryRun bool) {
	m.dryRun = dryRun
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
	finalReport := NewManagerReporter(cbpInfo)
	currentCBWP := fmt.Sprintf("当前共享带宽实例: %s", cbpInfo.ID)
	currentCBWPIn := fmt.Sprintf("平均流入带宽: %v Mbps;", rxDataPoint.Value)
	finalReport.AddContent(currentCBWPIn)
	currentCBWPOut := fmt.Sprintf("平均流出带宽: %v Mbps;", txDataPoint.Value)
	finalReport.AddContent(currentCBWPOut)
	log.Info().Msg(currentCBWP + currentCBWPIn + currentCBWPOut)
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
	var conclusion string
	// 当前共享带宽最大带宽速率,单位是Mbps
	currentMaxBandwidthRate := math.Max(rxDataPoint.Value, txDataPoint.Value)
	if currentMaxBandwidthRate > float64(cbpInfo.MaxBandwidth) {
		conclusion = "带宽高峰，需要缩容"
		log.Warn().Msg(conclusion)
		finalReport.AddConclusion(conclusion)
		err := m.ScaleDown(currentMaxBandwidthRate, finalReport)
		if err != nil {
			log.Err(err)
		}
		return
	}
	// 如果MinBandwidth是30，MaxBandwidth是70的话，小于（30+70）/2=50Mbps就会触发扩容
	if currentMaxBandwidthRate < float64((cbpInfo.MinBandwidth+cbpInfo.MaxBandwidth)/2) {
		conclusion = "带宽低谷，需要扩容"
		log.Warn().Msg(conclusion)
		finalReport.AddConclusion(conclusion)
		err := m.ScaleUp(currentMaxBandwidthRate, finalReport)
		if err != nil {
			log.Err(err)
		}
		return
	}
	//无需扩容,也无需缩容
	conclusion = "无需扩容,也无需缩容"
	log.Info().Msg(conclusion)
	if e := log.Debug(); e.Enabled() {
		if len(m.dingtalkNotifyToken) > 0 {
			finalReport.AddConclusion(conclusion)
			finalReport.ExportToDingTalk(m.dingtalkNotifyToken)
		}
	}
}

// ScaleUp 扩容:将低带宽EIP加入共享带宽
func (m *Manager) ScaleUp(currentBandwidthRate float64, reporter *ManagerReporter) (err error) {
	cbpInfo := m.cbp
	cbwpID := cbpInfo.ID
	//获取当前 region 下未绑定共享带宽的IP列表
	currentUnbindEIPs, err := m.sdk.GetCurrentEipAddressesExceptCBWP(cbwpID)
	if err != nil {
		return err
	}
	if len(currentUnbindEIPs) < 1 {
		return fmt.Errorf("len(currentUnbindEIPs) == 0")
	}
	log.Info().Msgf("len(currentUnbindEIPs):%v;currentUnbindEIPs:%v", len(currentUnbindEIPs), currentUnbindEIPs)
	for k, v := range currentUnbindEIPs {
		log.Info().Msgf("currentUnbindEIPs[%v] ;EIP: %s ;EIPID: %s", k, v.IpAddress, v.AllocationId)
	}
	var eipWaitLock sync.WaitGroup
	checkFrequency := cbpInfo.CheckFrequency
	eipWaitLock.Add(len(currentUnbindEIPs))
	var eipAvgList []model.EipAvgBandwidthInfo
	for _, eipInfo := range currentUnbindEIPs {
		go func(eip vpc.EipAddress, wg *sync.WaitGroup) {
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
		}(eipInfo, &eipWaitLock)
	}
	eipWaitLock.Wait()
	log.Info().Msgf("eipAvgList:%v", eipAvgList)
	//根据剩余带宽动态规划
	bandwidthLimit := (cbpInfo.MinBandwidth+cbpInfo.MaxBandwidth)/2 - int(currentBandwidthRate)
	currentSituation := fmt.Sprintf("剩余可用带宽 bandwidthLimit: %v Mbps", bandwidthLimit)
	log.Info().Msgf(currentSituation)
	reporter.AddContent(currentSituation)
	bestPublicIpAddress, err := model.NewBestPublicIpAddress(bandwidthLimit, eipAvgList)
	if err != nil {
		return err
	}
	// 动态优化求最优IP
	bestEIPs := bestPublicIpAddress.FindBestWithoutBrain()
	if len(bestEIPs) < 1 {
		msg := "剩余带宽不够绑定新的EIP"
		log.Info().Msg(msg)
		if len(m.dingtalkNotifyToken) > 0 {
			reporter.AddConclusion(msg)
			reporter.ExportToDingTalk(m.dingtalkNotifyToken)
		}
		return nil
	}
	if !m.dryRun {
		for _, eipInfo := range bestEIPs {
			reporter.AddContent(fmt.Sprintf(ADD_EIP_TEMPLATE, eipInfo.AllocationId, eipInfo.IpAddress))
			m.sdk.AddCommonBandwidthPackageIp(cbpInfo.ID, eipInfo.AllocationId)
		}
	}
	if len(m.dingtalkNotifyToken) > 0 {
		reporter.ExportToDingTalk(m.dingtalkNotifyToken)
	}
	return nil
}

//ScaleDown 缩容:将高带宽EIP移除共享带宽
func (m *Manager) ScaleDown(currentBandwidthRate float64, reporter *ManagerReporter) (err error) {
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
		go func(eip vpc.PublicIpAddresse, wg *sync.WaitGroup) {
			defer wg.Done()
			avgBandwidth, err := m.sdk.DescribeEipAvgMonitorData(eip.AllocationId, checkFrequency)
			//FIXME: 局部失败要怎么处理
			if err != nil {
				log.Err(err)
				return
			}
			eIPBandwidthInfo := model.EipAvgBandwidthInfo{
				IpAddress:    eip.IpAddress,
				AllocationId: eip.AllocationId,
				Value:        avgBandwidth,
			}
			log.Info().Msgf("eIPBandwidthInfo: IpAddress:%s ;AllocationId:%s ;avgBandwidth(Mbps):%v",
				eIPBandwidthInfo.IpAddress,
				eIPBandwidthInfo.AllocationId,
				eIPBandwidthInfo.Value)
			eipAvgList = append(eipAvgList, eIPBandwidthInfo)
		}(eipInfo, &eipWaitLock)
	}
	eipWaitLock.Wait()
	log.Info().Msgf("len(eipAvgList):%v", len(eipAvgList))
	bestPublicIpAddress, err := model.NewBestPublicIpAddress(cbpInfo.MinBandwidth, eipAvgList)
	if err != nil {
		return err
	}
	var conclusion string
	//进行动态优化
	bestIPs := bestPublicIpAddress.FindBestWithoutBrain()
	if len(bestIPs) < 1 {
		if len(eipAvgList) > 0 {
			conclusion = "结论：你这个程序有 bug 了"
			log.Warn().Msg(conclusion)
		} else {
			conclusion = "结论：没啥好优化的,再见"
			log.Info().Msg(conclusion)
		}
		reporter.AddStep(conclusion)
		reporter.ExportToDingTalk(m.dingtalkNotifyToken)
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
			reporter.AddContent(fmt.Sprintf(REMOVE_EIP_TEMPLATE, entity.AllocationId, entity.IpAddress))
			lowestEIPsAddress = append(lowestEIPsAddress, ip.AllocationId)
		}
	}
	log.Debug().Msgf("lowestEIPs:%v", lowestEIPs)
	if len(m.dingtalkNotifyToken) > 0 {
		reporter.ExportToDingTalk(m.dingtalkNotifyToken)
	}
	if !m.dryRun {
		m.sdk.RemoveCommonBandwidthPackageIps(cbpInfo.ID, lowestEIPsAddress)
	}
	return nil
}
