package alerter

import (
	"github.com/imroc/req/v3"
	"github.com/xmapst/logx"
)

var _ IAlert = (*sBase)(nil)

type sBase struct {
	http *req.Client
	url  string
	ak   string
	sk   string
}

func (s *sBase) Send(text string) {
	logx.Infoln(text)
}

func (s *sBase) SetUrl(url string) {
	s.url = url
}

func (s *sBase) SetAk(ak string) {
	s.ak = ak
}

func (s *sBase) SetSk(sk string) {
	s.sk = sk
}
