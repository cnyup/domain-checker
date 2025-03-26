package core

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/xmapst/logx"

	"github.com/busybox-org/cert-checker/internal/osext"
	"github.com/busybox-org/cert-checker/internal/selfupdate"
)

func (p *sProgram) selfUpdate() {
	_, err := p.cron.AddFunc("@every 30m", func() {
		p.doUpdate()
	})
	if err != nil {
		logx.Warnln(err)
		return
	}
}

func (p *sProgram) doUpdate() {
	if p.sURL == "" {
		logx.Debugln("automatic updates are not turned on")
		return
	}

	if p.sHash == nil {
		p.localSha256sum()
	}

	// 获取远端hash, example: https://oss.yfdou.com/tools/AutoExecFlow/latest.sha256sum
	checksum, name := p.remoteSha256sum(fmt.Sprintf("%s/latest.sha256sum", p.sURL))
	if checksum == nil || p.sHash == nil || bytes.Equal(p.sHash, checksum) {
		// 获取不到本地, 远端hash, 或者两者相同不更新
		return
	}

	opts := selfupdate.Options{
		Checksum: checksum,
	}
	// 检查权限
	if err := opts.CheckPermissions(); err != nil {
		logx.Warnln(err)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%s/%s", p.sURL, name))
	if err != nil {
		logx.Warnln(err)
		return
	}
	if resp == nil {
		logx.Warnln("resp is nil")
		return
	}
	defer func(Body io.ReadCloser) {
		err = Body.Close()
		if err != nil {
			logx.Errorln(err)
		}
	}(resp.Body)

	err = selfupdate.PrepareAndCheckBinary(resp.Body, opts)
	if err != nil {
		logx.Warnln(err)
		return
	}

	// 更新文件
	err = selfupdate.CommitBinary(opts)
	if err != nil {
		logx.Errorln(err)
		return
	}

	// 依赖systemd或者service守护
	os.Exit(200)
}

func (p *sProgram) localSha256sum() {
	path, err := osext.Executable()
	if err != nil {
		logx.Warnln(err)
		return
	}
	hash := sha256.New()
	file, err := os.Open(path)
	if err != nil {
		logx.Warnln(err)
		return
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			logx.Errorln(err)
		}
	}(file)
	_, err = io.Copy(hash, file)
	if err != nil {
		logx.Warnln(err)
		return
	}
	p.sHash = hash.Sum([]byte{})
	return
}

func (p *sProgram) remoteSha256sum(url string) ([]byte, string) {
	resp, err := http.Get(url)
	if resp != nil {
		defer func(Body io.ReadCloser) {
			err = Body.Close()
			if err != nil {
				logx.Errorln(err)
			}
		}(resp.Body)
	}
	if err != nil {
		logx.Warnln(err)
		return nil, ""
	}
	body := bufio.NewReader(resp.Body)
	for {
		var line []byte
		line, _, err = body.ReadLine()
		if err == io.EOF {
			return nil, ""
		}
		if err != nil {
			logx.Warnln(err)
		}
		if bytes.Contains(line, []byte(runtime.GOOS)) &&
			bytes.Contains(line, []byte(runtime.GOARCH)) {
			fields := bytes.Fields(line)
			if len(fields) != 2 {
				continue
			}
			src := fields[0]
			var n int
			n, err = hex.Decode(src, src)
			if err == nil {
				return src[:n], filepath.Base(string(fields[1]))
			} else {
				logx.Warnln(err)
			}
		}
	}
}
