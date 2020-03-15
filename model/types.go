package model

type CommonBandwidthPackage struct {
	// ID 共享带宽ID
	ID string `json:"id"`
	// MaxBandwidth 可接受最大带宽(用于缩容)
	MaxBandwidth int `yaml:"maxBandwidth"`
	// MinBandwidth 可接受最小带宽(用于扩容)
	MinBandwidth   int    `yaml:"minBandwidth"`
	CheckFrequency string `yaml:"checkFrequency"`
	Region         string `yaml:"region"`
}

type AliyunConfig struct {
	Region       string `yaml:"region"`
	AccessKeyId  string `yaml:"accessKeyId"`
	AccessSecret string `yaml:"accessSecret"`
}

type Datapoint struct {
	Timestamp  int64   `json:"timestamp"`
	UserId     string  `json:"userId"`
	InstanceId string  `json:"instanceId"`
	Value      float64 `json:"Value"`
	// just for copy
	// `json:""`
}
