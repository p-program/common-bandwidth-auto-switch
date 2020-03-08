package model

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadYAML(t *testing.T) {
	config := NewProjectConfig()
	path := path.Join("../", "config.yaml")
	err := config.LoadYAML(path)
	assert.Nil(t, err)
	aliyunConfig := config.AliyunConfig
	assert.Greater(t, len(aliyunConfig.Region), 0)
	assert.Greater(t, len(aliyunConfig.AccessKeyId), 0)
	assert.Greater(t, len(aliyunConfig.AccessSecret), 0)
}

func TestRemoveYAMLcomment1(t *testing.T) {
	text := `dklfjalsdjflskdjf WAIT A NEW LINE
###### bilibili
//// dilidili
alibalalabalia`
	t.Log(text)
	result := removeYAMLcomment(text)
	t.Log(result)
}

func TestRemoveYAMLcomment2(t *testing.T) {
	text := `dklfjalsdjflskdjf WAIT A NEW LINE
alibalalabalia`
	t.Log(text)
	result := removeYAMLcomment(text)
	t.Log(result)
}
