package manager

import (
	"errors"
	"math"
	"sync"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/vpc"
	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
)

// Manager 控制终端
type Manager struct {
	// sdk *model.AliyunConfig
	sdk *aliyun.AliyunSDK
	cbp *model.CommonBandwidthPackage
}

func NewManager(sdk *aliyun.AliyunSDK, cbp *model.CommonBandwidthPackage) *Manager {
	return &Manager{
		sdk: sdk,
		cbp: cbp,
	}
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
	}(&currentRateWG)
	go func(w *sync.WaitGroup) {
		txDataPoint, err2 = m.sdk.GetAvgTxRate(m.cbp)
	}(&currentRateWG)
	currentRateWG.Wait()
	// 共享带宽信息
	cbpInfo := m.cbp
	log.Info().Msgf("当前共享带宽实例: %s ;平均流入带宽: %d ;平均流出带宽: %d ;",
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
	if currentMaxBandwidthRate < float64(cbpInfo.MinBandwidth) {
		//带宽低谷，需要扩容
		m.ScaleUp(currentMaxBandwidthRate)
		return
	}
	//无需扩容,也无需缩容
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
	bestPublicIpAddress, err := model.NewBestPublicIpAddress(m.cbp.MinBandwidth, eipAvgList)
	if err != nil {
		return err
	}
	//TODO
	return nil
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
	lowestEIPs := bestPublicIpAddress.FindLowestPublicIpAddress()
	var ipInstanceIds []string
	for _, eIP := range lowestEIPs {
		ipInstanceIds = append(ipInstanceIds, eIP.AllocationId)
	}
	m.sdk.RemoveCommonBandwidthPackageIps(m.cbp.ID, ipInstanceIds)
	return nil
}
