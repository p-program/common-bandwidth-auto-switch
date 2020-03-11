package model

// EipAvgBandwidthInfo PublicIpAddresse的简化版,为了方便排序
// 参考 https://segmentfault.com/a/1190000008062661
type EipAvgBandwidthInfo struct {
	IpAddress string `json:"IpAddress" `
	//AllocationId EIP ID
	AllocationId string `json:"AllocationId"`
	//Value 带宽值
	Value float64 `json:"Value"`
}

type EipAvgBandwidthInfos []EipAvgBandwidthInfo

// 获取此 slice 的长度
func (list EipAvgBandwidthInfos) Len() int { return len(list) }

// 根据带宽升序排序 （此处按照自己的业务逻辑写）
func (list EipAvgBandwidthInfos) Less(i, j int) bool {
	return list[j].Value > list[i].Value
}

// 交换数据
func (list EipAvgBandwidthInfos) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

// PickSomeEIP 根据允许的最大
func (list *EipAvgBandwidthInfos) PickSomeEIP(maxBandwidth int) *EipAvgBandwidthInfos {

	return list
}
