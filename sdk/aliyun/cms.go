package aliyun

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/aliyun/alibaba-cloud-sdk-go/services/cms"
	"github.com/rs/zerolog/log"
	"github.com/zeusro/common-bandwidth-auto-switch/model"
)

func (sdk *AliyunSDK) GetCMSClient() *cms.Client {
	client, err := cms.NewClientWithAccessKey(sdk.config.Region, sdk.config.AccessKeyId, sdk.config.AccessSecret)
	if err != nil {
		log.Err(err)
		return nil
	}
	return client
}

const (
	format1 = "2006-01-02 15:04:05"
)

// DescribeMetricList see
//
// 共享带宽监控项 https://help.aliyun.com/document_detail/28619.html#title-hjj-o69-elv
func (sdk *AliyunSDK) DescribeMetricList(metricName string, cbp *model.CommonBandwidthPackage) (dataPoints []model.Datapoint, err error) {
	now := time.Now()
	frequency, err := time.ParseDuration(cbp.CheckFrequency)
	log.Info().Msgf("duration: %s", frequency.String())
	if err != nil {
		return nil, err
	}
	startTime := now.Add(-frequency).Unix() * 1000
	client := sdk.GetCMSClient()
	request := cms.CreateDescribeMetricListRequest()
	request.Scheme = "https"
	request.Namespace = "acs_bandwidth_package"
	request.MetricName = metricName
	request.Dimensions = fmt.Sprintf(`[{"instanceId": "%s"}]`, cbp.ID)
	request.StartTime = strconv.FormatInt(startTime, 10)
	log.Info().Msgf("request.StartTime: %s", request.StartTime)
	request.EndTime = strconv.FormatInt(now.Unix()*1000, 10)
	log.Info().Msgf("request.EndTime: %s", request.EndTime)
	// 60 秒一个周期
	request.Period = "60"
	response, err := client.DescribeMetricList(request)
	if err != nil {
		return dataPoints, err
	}
	if !response.Success {
		err = errors.New(response.Message)
		return dataPoints, err
	}
	// log.Info().Msgf("response: %v", response)
	// log.Info().Msgf("response.Datapoints: %s", response.Datapoints)
	datas := []byte(response.Datapoints)
	err = json.Unmarshal(datas, &dataPoints)
	return
}

func getAvgDatapoints(dataPoints []model.Datapoint) (*model.Datapoint, error) {
	dataPointsLen := len(dataPoints)
	log.Info().Msgf("dataPointsLen:%v", dataPointsLen)
	if dataPointsLen < 1 {
		err := errors.New("len(dataPoints) == 0")
		return nil, err
	}
	var sum float64 = 0
	for _, v := range dataPoints {
		sum += v.Value
	}
	result := &model.Datapoint{
		InstanceId: dataPoints[0].InstanceId,
		UserId:     dataPoints[0].UserId,
	}
	// 1 mbps = 1048576 bps
	result.Value = sum / float64(dataPointsLen*1048576)
	log.Info().Msgf("共享带宽实例ID:%s ;平均带宽:%v", result.InstanceId, result.Value)
	return result, nil
}

// GetAvgRxRate 共享带宽 流入带宽
func (sdk *AliyunSDK) GetAvgRxRate(cbp *model.CommonBandwidthPackage) (*model.Datapoint, error) {
	dataPoints, err := sdk.DescribeMetricList("net_rx.rate", cbp)
	if err != nil {
		return nil, err
	}
	return getAvgDatapoints(dataPoints)
}

// GetAvgTxRate 共享带宽 流出带宽
func (sdk *AliyunSDK) GetAvgTxRate(cbp *model.CommonBandwidthPackage) (*model.Datapoint, error) {
	dataPoints, err := sdk.DescribeMetricList("net_tx.rate", cbp)
	if err != nil {
		return nil, err
	}
	return getAvgDatapoints(dataPoints)
}
