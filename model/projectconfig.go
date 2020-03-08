package model

import (
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v2"
)

type ProjectConfig struct {
	CommonBandwidthPackages []CommonBandwidthPackage `yaml:"commonBandwidthPackages"`
	AliyunConfig            AliyunConfig             `yaml:"aliyun"`
	// for copy
	// `yaml:""`
}

func NewProjectConfig() *ProjectConfig {
	return &ProjectConfig{}

}

func (config *ProjectConfig) LoadYAML(path string) error {
	content, err := ioutil.ReadFile(path)
	contentWithoutComment := []byte(removeYAMLcomment(string(content)))
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(contentWithoutComment, &config)
	if err != nil {
		return err
	}
	return err
}

func removeYAMLcomment(oldText string) (resultText string) {
	array := strings.Split(oldText, LineBreak)
	// https://medium.com/@thuc/8-notes-about-strings-builder-in-golang-65260daae6e9
	sb := &strings.Builder{}
	for _, singleLine := range array {
		if strings.HasPrefix(singleLine, "#") || strings.HasPrefix(singleLine, "//") {
			continue
		}
		sb.WriteString(singleLine + LineBreak)
	}
	resultText = sb.String()
	return resultText
}
