package checker

import (
	"crypto/x509"
	"encoding/pem"
	"errors"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cnyup/domain-checker/utils"
)

type Checker struct {
	EcsInfo         *ecsInfo      `json:"ecs_info"`
	ExpireDomain    []*domainInfo `json:"expire_domain"`
	ThresholdDomain []*domainInfo `json:"threshold_domain"`

	Url string
}

type domainInfo struct {
	DomainName  string `json:"domain_name"`
	ExpiredDays int    `json:"expired_days"`
}

type ecsInfo struct {
	Name  string `json:"name"`
	LanIp string `json:"lan_ip"`
	WanIp string `json:"wan_ip"`
}

func NewChecker() *Checker {
	c := &Checker{
		EcsInfo:         &ecsInfo{},
		ExpireDomain:    make([]*domainInfo, 0),
		ThresholdDomain: make([]*domainInfo, 0),
	}
	c.init()
	return c
}

func (c *Checker) init() {
	var err error
	// TODO 耗时很长 需要优化
	// ecs 信息初始化
	if c.EcsInfo.LanIp, err = utils.GetInternalIP(); err != nil {
		slog.Warn("获取内网ip失败")
	}
	if c.EcsInfo.WanIp, _ = utils.GetExternalIP(); err != nil {
		slog.Warn("获取外网ip失败")
	}
	if c.EcsInfo.Name, err = os.Hostname(); err != nil {
		slog.Warn("获取主机名失败")
	}
}

// checkCert
//
//	@Description: 证书文件检测
//	@receiver c
//	@param file
//	@param days
//	@return error
func (c *Checker) checkCert(file string, days int) error {
	pemTmp, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	// decode
	certBlock, _ := pem.Decode(pemTmp)
	if certBlock == nil {
		return errors.New("decode cert file failed")

	}
	certBody, err := x509.ParseCertificate(certBlock.Bytes)
	if err != nil {
		return err
	}

	var info domainInfo
	// 获取域名信息
	info.DomainName = certBody.DNSNames[0]
	// 获取时间
	info.ExpiredDays = int(certBody.NotAfter.Sub(time.Now()).Hours() / 24)

	switch {
	case info.ExpiredDays < 0:
		// 已经过期
		c.ExpireDomain = append(c.ExpireDomain, &info)
	case info.ExpiredDays > 0 && info.ExpiredDays < days:
		// 触发阈值
		c.ThresholdDomain = append(c.ThresholdDomain, &info)
	default:
		return nil
	}

	return nil
}

// CheckByFile
//
//	@Description: 单文件检测 结果通过回调函数处理
//	@receiver c
//	@param file
//	@param days
//	@param fn
//	@return error
func (c *Checker) CheckByFile(file string, days int, fn func()) error {
	err := c.checkCert(file, days)
	if err != nil {
		return err
	}

	fn()
	return nil
}

// CheckByDir
//
//	@Description: 目录检测 结果通过回调函数处理
//	@receiver c
//	@param dir
//	@param su
//	@param days
//	@param fn
//	@return error
func (c *Checker) CheckByDir(dir, su string, days int, fn func()) error {
	// 忽略大小写
	su = strings.ToLower(su)

	// 检索目录
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		// 排除其他后缀文件

		if !strings.HasSuffix(path, su) {
			return nil
		}
		// 对crt文件进行处理
		err = c.checkCert(path, days)
		if err != nil {
			return err
		}
		return nil
	})

	fn()
	return nil
}
