package checker

import (
	"bytes"
	"html/template"
	"io"
	"log/slog"
	"net/http"
)

// Check2DingTalkByFile
//
//	@Description:  单文件检测 发送到钉钉
//	@receiver c
//	@param file
//	@param days
//	@return error
func (c *Checker) Check2DingTalkByFile(file string, days int) error {
	err := c.CheckByFile(file, days, c.dingTalk)
	if err != nil {
		return err
	}
	return nil
}

// Check2DingTalkByDir
//
//	@Description: 目录检测 发送到钉钉
//	@receiver c
//	@param dir
//	@param su
//	@param days
//	@return error
func (c *Checker) Check2DingTalkByDir(dir, su string, days int) error {
	err := c.CheckByDir(dir, su, days, c.dingTalk)
	if err != nil {
		return err
	}
	return nil
}

// dingTalk
//
//	@Description: 生成钉钉消息体 发送钉钉消息
//	@receiver c
func (c *Checker) dingTalk() {
	tmpl, _ := template.New("").Parse(DingTemplate)

	var buf bytes.Buffer

	defer buf.Reset()
	err := tmpl.Execute(&buf, c)
	if err != nil {
		slog.Error("template error:", err)
		return
	}
	reader := bytes.NewReader(buf.Bytes())
	send(c.Url, reader)
}

func send(url string, data *bytes.Reader) {
	req, _ := http.NewRequest("POST", url, data)

	req.Header.Add("content-type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			return
		}
	}(res.Body)
	body, _ := io.ReadAll(res.Body)

	slog.Info(string(body))
}
