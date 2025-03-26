package core

import (
	"bytes"
	"os"

	"github.com/kardianos/service"
	"github.com/robfig/cron/v3"
	"github.com/spf13/pflag"
	"github.com/xmapst/logx"

	"github.com/busybox-org/cert-checker/internal/alerter"
	"github.com/busybox-org/cert-checker/internal/core/checker"
	"github.com/busybox-org/cert-checker/internal/resolvers"
)

type sProgram struct {
	flags *pflag.FlagSet
	cron  *cron.Cron
	alert alerter.IAlert
	check checker.IChecker
	sHash []byte
	sURL  string
	// ecs info
	hostname string
	lanIP    string
	wanIP    string
}

func New(flags *pflag.FlagSet) service.Interface {
	daemon := &sProgram{
		flags: flags,
		alert: alerter.New(flags.Lookup("alert_type").Value.String()),
	}
	daemon.init()
	return daemon
}
func (p *sProgram) init() {
	p.cron = cron.New(cron.WithParser(cron.NewParser(
		cron.SecondOptional | cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor,
	)))
	p.cron.Start()
	p.selfUpdate()
	var err error
	// 获取主机名
	p.hostname, err = os.Hostname()
	if err != nil {
		logx.Warnf("获取主机名失败: %v", err)
		p.hostname = "unknown"
	}
	// 获取内网IP
	p.lanIP, err = resolvers.GetInternalIP()
	if err != nil {
		logx.Warnf("获取内网ip失败: %v", err)
		p.lanIP = "unknown"
	}
	p.wanIP, err = resolvers.GetExternalIP()
	if err != nil {
		logx.Warnf("获取外网ip失败: %v", err)
		p.wanIP = "unknown"
	}
	logx.Debugf("hostname: %s, lan_ip: %s, wan_ip: %s", p.hostname, p.lanIP, p.wanIP)
	alertType := p.flags.Lookup("alert_type").Value.String()
	p.alert = alerter.New(alertType)
	alertAk := p.flags.Lookup("alert_ak").Value.String()
	p.alert.SetAk(alertAk)
	alertSK := p.flags.Lookup("alert_sk").Value.String()
	p.alert.SetSk(alertSK)
}

func (p *sProgram) Start(service.Service) error {
	paths, err := p.flags.GetStringSlice("path")
	if err != nil {
		return err
	}
	suffix := p.flags.Lookup("suffix").Value.String()
	days, err := p.flags.GetInt("days")
	if err != nil {
		return err
	}
	p.check = checker.New(suffix)
	spec := p.flags.Lookup("cron").Value.String()
	_, err = p.cron.AddFunc(spec, func() {
		logx.Infof("开始检查证书...")
		var res []*checker.Response
		res, err = p.check.CheckCerts(paths...)
		if err != nil {
			logx.Warnf("检查证书失败: %v", err)
			return
		}
		var data = map[string]any{
			"EcsInfo": map[string]any{
				"Name":  p.hostname,
				"LanIp": p.lanIP,
				"WanIp": p.wanIP,
			},
			"ExpireDomain":    []any{},
			"ThresholdDomain": []any{},
		}
		for _, v := range res {
			if v.ExpiredDays < 0 {
				data["ExpireDomain"] = append(data["ExpireDomain"].([]any), map[string]any{
					"DomainName":  v.DomainName,
					"ExpiredDays": v.ExpiredDays,
				})
				continue
			}
			if v.ExpiredDays <= days {
				data["ThresholdDomain"] = append(data["ThresholdDomain"].([]any), map[string]any{
					"DomainName":  v.DomainName,
					"ExpiredDays": v.ExpiredDays,
				})
			}
		}
		if len(data["ExpireDomain"].([]any)) <= 0 || len(data["ThresholdDomain"].([]any)) <= 0 {
			return
		}
		var buf bytes.Buffer
		defer buf.Reset()
		if err = tmpl.Execute(&buf, data); err != nil {
			logx.Errorln(err)
			return
		}
		p.alert.Send(buf.String())
		logx.Infof("证书检查完成...")
	})
	if err != nil {
		logx.Errorln(err)
		return err
	}
	return nil
}

func (p *sProgram) Stop(service.Service) error {
	p.cron.Stop()
	return nil
}
