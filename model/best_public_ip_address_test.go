package model

import (
	"os"
	"testing"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

func init() {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

/*
00 40 41 42 43 44 45 46 47 48 49
02 02 02 02 02 02 02 02 02 02 02
05 07 07 07 07 07 07 07 07 07 07
10 17 17 17 17 17 17 17 17 17 17
15 32 32 32 32 32 32 32 32 32 32
20 35 35 35 35 35 35 35 35 35 35
21 35 41 41 41 41 41 41 41 41 41
*/
func TestDynamic1(t *testing.T) {
	bandwidthInfos := prepareEipAvgBandwidthInfos1()
	findBestPolicy, err := NewBestPublicIpAddress(40, bandwidthInfos)
	assert.Nil(t, err)
	//最优解： 21 + 20
	best := findBestPolicy.FindBestWithoutBrain()
	t.Log(best)
	findBestPolicy.print()
}

/*
00 40 41 42 43 44 45 46 47 48 49
20 20 20 20 20 20 20 20 20 20 20
21 20 41 41 41 41 41 41 41 41 41
31 20 41 41 41 41 41 41 41 41 41
*/
func TestDynamic2(t *testing.T) {
	bandwidthInfos := prepareEipAvgBandwidthInfos2()
	findBestPolicy, err := NewBestPublicIpAddress(40, bandwidthInfos)
	assert.Nil(t, err)
	findBestPolicy.dynamic()
	findBestPolicy.print()
}

func TestFindBest(t *testing.T) {
	bandwidthInfos := prepareEipAvgBandwidthInfos2()
	bestIPs, err := NewBestPublicIpAddress(40, bandwidthInfos)
	assert.Nil(t, err)
	best := bestIPs.FindBestWithoutBrain()
	bestIPs.print()
	t.Logf("len(best):%v", len(best))
	for j := 0; j < len(best); j++ {
		t.Logf("%v \n", best[j])
	}
}

func prepareEipAvgBandwidthInfos1() []EipAvgBandwidthInfo {
	bandwidthInfos := []EipAvgBandwidthInfo{
		EipAvgBandwidthInfo{"1.1.1.1", "a", float64(20)},
		EipAvgBandwidthInfo{"1.1.1.2", "", float64(10)},
		EipAvgBandwidthInfo{"1.1.1.3", "", float64(15)},
		EipAvgBandwidthInfo{"1.1.1.4", "", float64(5)},
		EipAvgBandwidthInfo{"1.1.1.5", "", float64(2)},
		EipAvgBandwidthInfo{"1.1.1.6", "", float64(21)},
	}
	return bandwidthInfos
}

func prepareEipAvgBandwidthInfos2() []EipAvgBandwidthInfo {
	bandwidthInfos := []EipAvgBandwidthInfo{
		EipAvgBandwidthInfo{"1.1.1.1", "", float64(21)},
		EipAvgBandwidthInfo{"1.1.1.2", "", float64(20)},
		EipAvgBandwidthInfo{"1.1.1.3", "", float64(31)},
	}
	return bandwidthInfos
}
