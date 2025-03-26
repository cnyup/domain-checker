package checker

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"time"
)

type IChecker interface {
	CheckCerts(paths ...string) ([]*Response, error)
}

type sChecker struct {
	suffix string
	day    int
}

type Response struct {
	Path        string `json:"path"`
	ExpiredDays int    `json:"expired_days"`
	DomainName  string `json:"domain_name"`
}

func New(suffix string) IChecker {
	return &sChecker{
		suffix: suffix,
	}
}

func (c *sChecker) CheckCerts(paths ...string) ([]*Response, error) {
	var res []*Response
	for _, path := range paths {
		_res, err := c.checkCert(path)
		if err != nil {
			return nil, err
		}
		res = append(res, _res...)
	}
	return res, nil
}

func (c *sChecker) checkCert(path string) ([]*Response, error) {
	var res []*Response
	// 判断路径是否为文件夹
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		_res, err := c.checkCertByFile(path)
		if err != nil {
			return nil, err
		}
		res = append(res, _res)
		return res, nil
	}
	err = c.WalkPath(path, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if !strings.HasSuffix(path, c.suffix) {
			return nil
		}
		_res, err := c.checkCertByFile(path)
		if err != nil {
			return err
		}
		res = append(res, _res)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (c *sChecker) checkCertByFile(path string) (*Response, error) {
	// check file suffix
	if !strings.HasSuffix(path, c.suffix) {
		return nil, fmt.Errorf("file suffix error, %s", path)
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	block, _ := pem.Decode(content)
	if block == nil {
		return nil, fmt.Errorf("decode cert file failed, %s", path)
	}
	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return nil, err
	}
	if len(cert.DNSNames) == 0 {
		return nil, fmt.Errorf("cert file dns name is empty, %s", path)
	}
	return &Response{
		Path:        path,
		ExpiredDays: int(cert.NotAfter.Sub(time.Now()).Hours() / 24),
		DomainName:  cert.DNSNames[0],
	}, nil
}
