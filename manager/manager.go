package manager

import (
	"errors"
	"math"
	"sync"

	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/sdk/aliyun"
)

// Manager 控制终端
type Manager struct {
	// sdk *model.AliyunConfig
	sdk aliyun.AliyunSDK
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
	// 当前共享带宽最大带宽速率
	currentMaxBandwidthRate := math.Max(rxDataPoint.Value, txDataPoint.Value)
	if currentMaxBandwidthRate > float64(cbpInfo.MaxBandwidth) {
		m.ScaleDown(currentMaxBandwidthRate)
		return
	}
	if currentMaxBandwidthRate < float64(cbpInfo.MinBandwidth) {
		m.ScaleUp(currentMaxBandwidthRate)
		return
	}
	//无需扩容,也无需缩容
}

func pickSomeEIP() {

}

// ScaleUp 扩容:将低带宽EIP加入共享带宽
func (m *Manager) ScaleUp(currentBandwidthRate float64) {

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
	// maxBandwidth := cbpInfo.MaxBandwidth
	// var targetRemovedEips []vpc.PublicIpAddresse

	//为了尽可能地减少 goroutine 创建,防止阿里云API限流,这里使用串行查询,当带宽满足要求时即可退出查询

	// var wg sync.WaitGroup
	// wg.Add()
	// for _, v := range eipList {

	// }
	//选取高带宽的EIP,然后将他们移除

}
