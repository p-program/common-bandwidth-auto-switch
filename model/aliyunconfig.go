package model

type AliyunConfig struct {
	Region       string `yaml:"region"`
	AccessKeyId  string `yaml:"accessKeyId"`
	AccessSecret string `yaml:"accessSecret"`
}
