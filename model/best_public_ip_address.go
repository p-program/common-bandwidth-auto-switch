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
	// 带宽偏移量,默认值10
	// offset int
	//动态规划网格，网格的元素为当前局部带宽最优解
	cellsMesh [X][Y]float64
	// 带宽(类似于背包问题的物品容量)
	bandwidths []float64
}

const (
	X = 10
	// Y 最大支持200条EIP检测数据
	Y = 200
	// DEFAULT_OFFSET=10
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
		for j := 1; j < Y; j++ {
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
	currentMaxBandwidth := float64(j + m.minBandwidth)
	var getValue = func(currentBandwidth int) float64 {
		// 离最小带宽越近越好,所以以负数形式返回
		return float64(m.minBandwidth - currentBandwidth)
	}
	if m.origin[i].Value > currentMaxBandwidth {
		return m.cellsMesh[i-1][j]
	}
	// 可放进背包时候，计算放入当前商品后的最大价值
	currentValue := getValue(i+m.minBandwidth) + m.cellsMesh[i-1][j-m.bandwidths[i]]

	return 0
}
