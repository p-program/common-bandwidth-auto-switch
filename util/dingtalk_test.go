package util

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	TEST_Token = "a12cc16245bf"
)

func TestDingMarkdown(t *testing.T) {
	d := NewDingTalk(TEST_Token)
	b := NewMarkdownBuilder()
	b.AddH1("共享带宽动态规划")
	b.AddText("&*%……&￥#￥%#￥@#￥@*&*（")
	markdownText := b.BuilderText()
	err := d.DingMarkdown("共享带宽动态规划", markdownText)
	assert.Nil(t, err)
}
