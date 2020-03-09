package model

import "fmt"

// BestPublicIpAddress 应用动态规划,寻求最佳EIP列表
// https://blog.csdn.net/runbat/article/details/94016554
type BestPublicIpAddress struct {
	// 原始EIP检测数据
	origin []EipAvgBandwidthInfo

	//best 最佳EIP池
	best []EipAvgBandwidthInfo
	//minBandwidth 最小带宽
	minBandwidth int
	//动态规划网格，网格的元素为当前局部带宽最优解
	cellsMesh [X][Y]float64
	// x轴是带宽(类似于背包问题的容量)
	bandwidths []float64
	// y轴是EIP检测数据的带宽
	eipAvgBandwidthInfos []float64
}

const (
	X = 10
	// Y 最大支持200条EIP检测数据
	Y = 200
)

func NewBestPublicIpAddress(minBandwidth int, maxBandwidth int) *BestPublicIpAddress {
	return &BestPublicIpAddress{
		minBandwidth: minBandwidth,
	}
}

func (m *BestPublicIpAddress) FindBestPublicIpAddress() {

}

func (m *BestPublicIpAddress) dynamic() {
	listLen := len(m.origin)

	// 初始化动态规划网格
	for i := 1; i < listLen; i++ {
		for j := m.minBandwidth; j < maxBandwidth; j++ {
			m.cellsMesh[i][j] = m.maxValue(i, j)
		}
	}
	for i := 0; i < Y; i++ {
		fmt.Printf("%v \n", m.cellsMesh[i])
	}

}

// 局部最优解
func (m *BestPublicIpAddress) maxValue(i, j int) float64 {
	// 当前商品无法放入背包，返回当前背包所能容纳的最大价值
	maxBandwidth := (m.minBandwidth + 10)
	if m.origin[i].Value > float64(j) {
		return m.cellsMesh[i-1][j]
	}
	return 0
}
