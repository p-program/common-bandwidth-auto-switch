package manager

import (
	"fmt"

	"github.com/zeusro/common-bandwidth-auto-switch/model"
	"github.com/zeusro/common-bandwidth-auto-switch/util"
)

type ManagerReporter struct {
	cbwpInfo     *model.CommonBandwidthPackage
	innerBuilder *util.MarkdownBuilder
}

func NewManagerReporter(cbwpInfo *model.CommonBandwidthPackage) *ManagerReporter {
	innerBuilder := util.NewMarkdownBuilder()
	r := &ManagerReporter{
		innerBuilder: innerBuilder,
		cbwpInfo:     cbwpInfo,
	}
	r.AddTitleLink()
	return r
}

func (r *ManagerReporter) ExportToDingTalk(token string) {
	ding := util.NewDingTalk(token)
	title := "共享带宽优化"
	ding.DingMarkdown(title, r.innerBuilder.BuilderText())
}

func (r *ManagerReporter) AddTitleLink() {
	r.innerBuilder.AddText("当前共享带宽实例:")
	u := fmt.Sprintf("https://vpcnext.console.aliyun.com/cbwp/%s/cbwps", r.cbwpInfo.Region)
	r.innerBuilder.AddLink(r.cbwpInfo.ID, u)
}

func (r *ManagerReporter) AddContent(content string) {
	r.innerBuilder.AddText(content)
}

func (r *ManagerReporter) AddConclusion(content string) {
	r.AddContent("结论：" + content)
}

func (r *ManagerReporter) AddStep(step string) {
	r.AddContent("举措：" + step)
}
