package model

import (
	"errors"
	"fmt"
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
	// 带宽偏移量,默认值10
	// offset int
	//动态规划网格，网格的元素为当前局部带宽最优解
	cellsMesh [Y][X]float64
	// 带宽(类似于背包问题的物品容量)
	bandwidths []float64
}

const (
	X = 10
	// Y 最大EIP检测数据
	Y = 200
	// DEFAULT_OFFSET=10
)

func NewBestPublicIpAddress(minBandwidth int, bandwidthInfos []EipAvgBandwidthInfo) (*BestPublicIpAddress, error) {
	if len(bandwidthInfos) > Y {
		err := errors.New("超过了可解范围")
		return nil, err
	}
	if len(bandwidthInfos) < 1 {
		err := errors.New("别来捣乱，OK?")
		return nil, err
	}
	bestIPs := &BestPublicIpAddress{
		minBandwidth: minBandwidth,
		origin:       bandwidthInfos,
	}
	return bestIPs, nil
}

func (m *BestPublicIpAddress) FindBestPublicIpAddress() []EipAvgBandwidthInfo {
	m.dynamic()
	m.findBack()
	return m.FindBestPublicIpAddress()
}

// FindLowestPublicIpAddress 求FindBestPublicIpAddress的差集
func (m *BestPublicIpAddress) FindLowestPublicIpAddress() []EipAvgBandwidthInfo {
	// TODO
	return m.FindLowestPublicIpAddress()
}

func (m *BestPublicIpAddress) dynamic() {
	//FIXME
	listLen := len(m.origin)
	// 初始化动态规划网格
	for y := 1; y < listLen; y++ {
		for x := 1; x < X; x++ {
			m.cellsMesh[y][x] = m.maxValue(y, x)
		}
	}
	for y := 0; y < listLen; y++ {
		fmt.Printf("%v \n", m.cellsMesh[y])
	}
}

// 局部最优解
func (m *BestPublicIpAddress) maxValue(y, x int) float64 {
	//FIXME
	// 当前商品无法放入背包，返回当前背包所能容纳的最大价值
	currentMaxBandwidth := float64(x + m.minBandwidth)
	if m.origin[y].Value > currentMaxBandwidth {
		return m.cellsMesh[y-1][x]
	}
	// 可放进背包时候，计算放入当前商品后的最大价值
	// 每一个EIP的价值都是1
	// 当前价值= EIP价值 + 剩余空间的价值
	// 剩余空间的价值=cell[i-1][j-当前商品的重量])
	currentValue := 1 + m.cellsMesh[y-1][x-1]
	if currentValue >= m.cellsMesh[y-1][x] {
		return currentValue
	}
	return m.cellsMesh[y-1][x]
}

// 回溯选择的商品方法
func (m *BestPublicIpAddress) findBack() []EipAvgBandwidthInfo {
	//FIXME
	col := X - 1
	for i := Y - 1; i > 0; i-- {
		if m.cellsMesh[i][col] > m.cellsMesh[i-1][col] {
			// selected[i] = 1
			m.best = append(m.best, m.origin[i])
			col = col - 1
		}
	}
	return m.best
}
