package model

// SimplePublicIpAddressInfo PublicIpAddresse的简化版,为了方便排序
// 参考 https://segmentfault.com/a/1190000008062661
type SimplePublicIpAddressInfo struct {
	IpAddress    string  `json:"IpAddress" `
	AllocationId string  `json:"AllocationId"`
	Value        float64 `json:"Value"`
}
type SimplePublicIpAddressInfos []SimplePublicIpAddressInfo

// 获取此 slice 的长度
func (list SimplePublicIpAddressInfos) Len() int { return len(list) }

// 根据带宽降序排序 （此处按照自己的业务逻辑写）
func (list SimplePublicIpAddressInfos) Less(i, j int) bool {
	return list[i].Value > list[j].Value
}

// 交换数据
func (list SimplePublicIpAddressInfos) Swap(i, j int) { list[i], list[j] = list[j], list[i] }

// PickSomeEIP 根据允许的最大
func (list *SimplePublicIpAddressInfos) PickSomeEIP(maxBandwidth int) *SimplePublicIpAddressInfos {

	return list
}
