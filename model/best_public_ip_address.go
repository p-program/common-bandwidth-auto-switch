package model

import (
	"errors"
	"fmt"

	"github.com/rs/zerolog/log"
)

// BestPublicIpAddress 应用动态规划,寻求最佳EIP列表
// https://blog.csdn.net/runbat/article/details/94016554
type BestPublicIpAddress struct {
	// 原始EIP检测数据
	origin []EipAvgBandwidthInfo
	//best 最佳EIP池
	best []EipAvgBandwidthInfo
	//minBandwidth 最小带宽
	minBandwidth int
	maxBandwidth int
	// cellsMesh 动态规划网格，网格的元素为当前局部带宽最优解（ XX Mbps）,由于 golang 不支持动态数组，这里只能初始化一个稍微大一点的数值
	cellsMesh [MAX_EIP][COL]float64
	//eipsLen 由于 golang 不支持动态数组,这里要确定二维边界
	eipsLen int
}

const (
	// COL 一维边界
	COL = 10
	// MAX_EIP 二维边界，最大EIP检测数据，尽量不要用这个初始化
	MAX_EIP = 102
)

func NewBestPublicIpAddress(minBandwidth int, bandwidthInfos []EipAvgBandwidthInfo) (*BestPublicIpAddress, error) {
	eipsLen := len(bandwidthInfos)
	if eipsLen > MAX_EIP {
		err := errors.New("IP列表超过了可解范围")
		return nil, err
	}
	if eipsLen < 1 {
		err := errors.New("别来捣乱，OK?")
		return nil, err
	}
	// 初始化动态规划网格
	var cells [MAX_EIP][COL]float64
	initBandwidth := float64(minBandwidth)
	for j := 1; j < COL; j++ {
		cells[0][j] = float64(initBandwidth)
		initBandwidth++
	}
	for i, v := range bandwidthInfos {
		//从第二行开始赋值
		cells[i+1][0] = float64(v.Value)
	}
	bestIPs := &BestPublicIpAddress{
		minBandwidth: minBandwidth,
		maxBandwidth: minBandwidth + COL,
		origin:       bandwidthInfos,
		eipsLen:      eipsLen,
		cellsMesh:    cells,
	}
	return bestIPs, nil
}

// FindBest 回溯选择,从最后一行最后一格开始推移
func (m *BestPublicIpAddress) FindBest() []EipAvgBandwidthInfo {
	m.dynamic()
	//TODO
	return m.best
}

func (m *BestPublicIpAddress) dynamic() {
	for i := 1; i < COL; i++ {
		for j := 1; j < m.eipsLen; j++ {
			log.Info().Msgf("i: %v ;j: %v", i, j)
			m.cellsMesh[i][j] = m.maxValue(i, j)
		}
	}
	for j := 0; j <= m.eipsLen; j++ {
		fmt.Printf("%v \n", m.cellsMesh[j])
	}
}

// 局部最优解
func (m *BestPublicIpAddress) maxValue(i, j int) float64 {
	// j 是带宽值
	lastColumnCell := m.cellsMesh[i-1][j]
	currentEIPBandwidth := m.origin[j].Value
	// 当前EIP超过带宽限制
	if currentEIPBandwidth > float64(j) {
		return lastColumnCell
	}
	// 剩余带宽=当前带宽上限-当前EIP的带宽
	// 剩余带宽足够容纳当前EIP+之前的IP
	remainingBandwidth := m.cellsMesh[0][j] - currentEIPBandwidth
	//从上一行中找到最符合需求的EIP
	currentCellBandwidth := currentEIPBandwidth
	//倒序遍历，最大值计入当前网格
	for k := len(m.cellsMesh[i-1]) - 1; k >= 0; k-- {
		//剩余带宽刚好能融入上一行的EIP的带宽
		if remainingBandwidth > m.cellsMesh[i-1][k] {
			currentCellBandwidth += m.cellsMesh[i-1][k]
			break
		}
	}
	if currentCellBandwidth >= lastColumnCell {
		return currentCellBandwidth
	}
	return lastColumnCell
}
