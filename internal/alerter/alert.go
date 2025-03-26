package alerter

import (
	"github.com/imroc/req/v3"
	"github.com/xmapst/logx"
)

type IAlert interface {
	SetUrl(url string)
	SetAk(ak string)
	SetSk(sk string)
	Send(text string)
}

func New(t string) IAlert {
	base := &sBase{
		http: req.NewClient().
			EnableHTTP3().
			EnableDumpAllAsync().
			ImpersonateChrome().
			SetLogger(logx.GetSubLogger()).
			SetCommonHeader("Content-Type", "application/json"),
	}
	switch t {
	case "dingtalk":
		return &sDingTalk{
			sBase: base,
		}
	default:
		return base
	}
}
