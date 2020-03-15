package util

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
)

// DingTalk https://ding-doc.dingtalk.com/doc#/serverapi2/qf2nxq
type DingTalk struct {
	token string
}

// DingTalkMsg struct
type DingTalkMsg struct {
	MsgType  string           `json:"msgtype"`
	Text     DingTalkText     `json:"text"`
	Markdown DingTalkMarkdown `json:"markdown"`
}

type DingTalkText struct {
	Content string `json:"content"`
}

type DingTalkMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

const (
	CONTENT_TYPE_JSON = "application/json"
)

func NewDingTalk(token string) *DingTalk {
	return &DingTalk{token}
}

func (d *DingTalk) DingMarkdown(title, markdownText string) (err error) {
	msg := &DingTalkMsg{
		MsgType: "markdown",
		Markdown: DingTalkMarkdown{
			Title: title,
			Text:  markdownText,
		},
	}
	return d.Ding(msg)
}

// Ding 钉钉消息推送
func (d *DingTalk) Ding(msg *DingTalkMsg) (err error) {
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		log.Warn().Err(fmt.Errorf("failed to marshal msg %v", msg))
		return err
	}
	b := bytes.NewBuffer(msgBytes)
	resp, err := http.Post(fmt.Sprintf("https://oapi.dingtalk.com/robot/send?access_token=%s", d.token), CONTENT_TYPE_JSON, b)
	if err != nil {
		//包装一下最终error
		finalErr := fmt.Errorf("failed to send msg to dingtalk. error: %s", err.Error())
		log.Err(finalErr)
		return err
	}
	defer resp.Body.Close()
	if resp != nil && resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("failed to send msg to dingtalk, because the response code is %d", resp.StatusCode)
		log.Err(err)
		return err
	}
	return nil
}
