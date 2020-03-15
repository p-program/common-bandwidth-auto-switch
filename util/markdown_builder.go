package util

import (
	"fmt"
	"strings"
)

type MarkdownBuilder struct {
	innerTextBuilder *strings.Builder
}

const (
	NEW_LINE = "\n\n"
	H1       = "# %s" + NEW_LINE
	H2       = "## %s" + NEW_LINE
	H3       = "### %s" + NEW_LINE
	H4       = "#### %s" + NEW_LINE
	H5       = "##### %s" + NEW_LINE
	BLOAD    = "**%s**" + NEW_LINE
	LINK     = "[%s](%s)" + NEW_LINE
)

func NewMarkdownBuilder() *MarkdownBuilder {
	return &MarkdownBuilder{
		innerTextBuilder: &strings.Builder{},
	}
}

func (b *MarkdownBuilder) BuilderText() string {
	return b.innerTextBuilder.String()
}

func (b *MarkdownBuilder) AddH1(h1 string) *MarkdownBuilder {
	b.innerTextBuilder.WriteString(fmt.Sprintf(H1, h1))
	return b
}
func (b *MarkdownBuilder) AddLink(title, u string) *MarkdownBuilder {
	b.innerTextBuilder.WriteString(fmt.Sprintf(LINK, title, u))
	return b
}

func (b *MarkdownBuilder) AddBload(content string) *MarkdownBuilder {
	b.innerTextBuilder.WriteString(fmt.Sprintf(BLOAD, content))
	return b
}

func (b *MarkdownBuilder) AddText(content string) *MarkdownBuilder {
	b.innerTextBuilder.WriteString(content + NEW_LINE)
	return b
}
