package alerter

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/xmapst/logx"
)

const dingtalkRobotUrl = "https://oapi.dingtalk.com/robot/send"

type sDingTalk struct {
	*sBase
}

func (d *sDingTalk) Send(text string) {
	if d.ak == "" {
		logx.Errorf("dingtalk ak is empty")
		return
	}
	if d.sBase.url == "" {
		d.http.SetBaseURL(dingtalkRobotUrl)
	}

	defer func() {
		d.http.CloseIdleConnections()
	}()

	req := d.http.NewRequest().SetQueryParam("access_token", d.ak)
	timestamp := time.Now().UnixNano() / 1e6
	req.SetQueryParam("timestamp", fmt.Sprintf("%d", timestamp))
	sign := d.getSign(timestamp)
	if sign != "" {
		req.SetQueryParam("sign", sign)
	}
	res, err := req.SetBody(map[string]any{
		"msgtype": "markdown",
		"markdown": map[string]any{
			"title": "域名证书即将过期",
			"text":  text,
		},
	}).Post("")
	if err != nil {
		logx.Errorf("dingtalk send error: %v", err)
		return
	}
	logx.Infoln(res.String())
}

func (d *sDingTalk) getSign(timestamp int64) (sign string) {
	if d.sk == "" {
		return
	}
	strToHash := fmt.Sprintf("%d\n%s", timestamp, d.sk)
	hmac256 := hmac.New(sha256.New, []byte(d.sk))
	hmac256.Write([]byte(strToHash))
	data := hmac256.Sum(nil)
	return base64.StdEncoding.EncodeToString(data)
}
