package model

import (
	"errors"
	"fmt"
	"sort"
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
	COL = 11
	// MAX_EIP 二维边界，最大EIP检测数据，尽量不要用这个初始化
	MAX_EIP = 102
)

func NewBestPublicIpAddress(minBandwidth int, bandwidthInfos EipAvgBandwidthInfos) (*BestPublicIpAddress, error) {
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
	cells := [MAX_EIP][COL]float64{}
	initBandwidth := float64(minBandwidth)
	for j := 1; j < COL; j++ {
		cells[0][j] = float64(initBandwidth)
		initBandwidth++
	}
	sort.Sort(bandwidthInfos)
	// fmt.Printf("len(cells): %v ;cap(cells): %v \n", len(cells[1]), cap(cells[1]))
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
	// fmt.Printf("len(cells): %v ;cap(cells): %v \n", len(bestIPs.cellsMesh[1]), cap(bestIPs.cellsMesh[1]))
	return bestIPs, nil
}

// FindBest 回溯选择,从最后一行最后一格开始推移
func (m *BestPublicIpAddress) FindBest() []EipAvgBandwidthInfo {
	m.dynamic()
	//TODO
	return m.best
}

func (m *BestPublicIpAddress) dynamic() {
	for i := 1; i <= m.eipsLen; i++ {
		for j := 1; j < COL; j++ {
			m.cellsMesh[i][j] = m.maxValue(i, j)
			// fmt.Printf("m.cellsMesh[%v][%v]: %v \n", i, j, m.cellsMesh[i][j])
		}
	}
}
func (m *BestPublicIpAddress) print() {
	for j := 0; j <= m.eipsLen; j++ {
		for _, v := range m.cellsMesh[j] {
			content := ""
			if v < float64(10) {
				content = "0"
			}
			// content = fmt.Sprintf("%s%v", content, v)
			fmt.Printf("%s%v ", content, v)
			// fmt.Printf("%v \n", m.cellsMesh[j])
		}
		fmt.Print("\n")
	}
}

// 局部最优解,只是近似最优解，不是最优解
func (m *BestPublicIpAddress) maxValue(i, j int) float64 {
	lastColumnCell := m.cellsMesh[i-1][j]
	currentEIPBandwidth := m.origin[i-1].Value
	bandwidthLimit := m.cellsMesh[0][j]
	// fmt.Printf("i: %v ; j: %v ;currentEIPBandwidth: %v ;bandwidthLimit: %v ;", i, j, currentEIPBandwidth, bandwidthLimit)
	// 当前EIP超过带宽限制
	if currentEIPBandwidth > bandwidthLimit {
		return lastColumnCell
	}
	if i == 1 {
		// 第2列要先特殊处理
		return currentEIPBandwidth
	}
	// 剩余带宽=当前带宽上限-当前EIP的带宽
	// 剩余带宽足够容纳当前EIP+之前的IP
	remainingBandwidth := bandwidthLimit - currentEIPBandwidth
	currentCellBandwidth := currentEIPBandwidth
	hasRemain := false
	//FIXME:从先前的元素中排列组合，求满足条件的最大值
	for k := COL - 1; k >= 0; k-- {
		//剩余带宽刚好能融入上一行的EIP的带宽
		if remainingBandwidth >= m.cellsMesh[i-1][k] {
			currentCellBandwidth += m.cellsMesh[i-1][k]
			hasRemain = true
			return currentCellBandwidth
		}
	}
	// 剩余带宽不够支持
	if !hasRemain {
		return lastColumnCell
	}
	return currentEIPBandwidth
}
