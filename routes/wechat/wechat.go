package wechat

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const wechatWekhookURL = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send"

func send(key string, content string) error {
	data, _ := json.Marshal(struct {
		MsgType  string   `json:"msgtype"`
		Markdown markdown `json:"markdown"`
	}{
		MsgType: "markdown",
		Markdown: markdown{
			Content: content,
		},
	})
	_, err := http.Post(fmt.Sprintf("%s?key=%s", wechatWekhookURL, key), "application/json", bytes.NewBuffer(data))
	return err
}
