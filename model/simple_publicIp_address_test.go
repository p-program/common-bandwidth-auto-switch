package model

import (
	"sort"
	"testing"
)

func TestSort(t *testing.T) {
	array := prepareSimplePublicIpAddressInfoArray()
	sort.Sort(array)
	//结果为降序,带宽高的显示靠前
	t.Logf("sorted array: %v", array)
}

func TestArray(t *testing.T) {
	info1 := SimplePublicIpAddressInfo{
		Value: float64(1),
	}
	info2 := SimplePublicIpAddressInfo{
		Value: float64(2),
	}
	list := []SimplePublicIpAddressInfo{info1, info2}
	item0 := list[0]
	item0.Value = float64(666666.66666)
	t.Logf("%v", item0.Value)
	t.Logf("%v", list[0].Value)

}
