package resolvers

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/busybox-org/cert-checker/internal/resolvers/base"
	"github.com/busybox-org/cert-checker/internal/resolvers/targets"
)

func GetInternalIP() (string, error) {
	// 获取本地网络接口
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}

	for _, iface := range ifaces {
		// 过滤掉未激活的接口
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}

		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.IsPrivate() {
				return ipnet.IP.String(), nil
			}
		}
	}
	return "", fmt.Errorf("未找到内网 IP 地址")
}

func GetExternalIP() (string, error) {
	res, err := optimizedRetrieve(3*time.Second, 0.6)
	if err != nil {
		return "", err
	}
	return res.IP.String(), nil
}

func optimizedRetrieve(timeout time.Duration, threshold float64) (*base.ScoredIP, error) {
	// 创建上下文，并在满足条件时取消其它操作
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// 预先计算总权重
	var sumOfWeight float64
	retrievables := targets.IPv4Retrievables()
	for _, sipr := range retrievables {
		sumOfWeight += sipr.Weight
	}

	// 创建用于传递 IP 检索结果的通道
	results := make(chan base.ScoredIPWithMaxScore)

	// 使用 WaitGroup 等待所有检索 goroutine 完成
	var wg sync.WaitGroup
	for _, sips := range retrievables {
		wg.Add(1)
		go func(sipr base.ScoredIPRetrievable) {
			defer wg.Done()
			sip, err := sipr.RetrieveIPWithScoring(ctx)
			if err != nil {
				return
			}
			// 将结果发送到通道（若 ctx 已取消，则不阻塞发送）
			select {
			case results <- *sip:
			case <-ctx.Done():
			}
		}(sips)
	}

	// 在所有任务完成后关闭 results 通道
	go func() {
		wg.Wait()
		close(results)
	}()
	type scoreAgg struct {
		score    float64
		maxScore float64
	}
	// 用 map 存储每个 IP 的累加得分和最大可能得分
	aggMap := make(map[string]*scoreAgg)
	// 遍历结果，并动态调整总权重
	for sip := range results {
		key := sip.IP.String()
		if _, exists := aggMap[key]; !exists {
			aggMap[key] = &scoreAgg{}
		}
		aggMap[key].score += sip.Score
		aggMap[key].maxScore += sip.MaxScore

		// 计算当前 IP 的得分率
		currentScore := aggMap[key].score / aggMap[key].maxScore
		if currentScore > threshold {
			cancel() // 提前终止其他任务
			return &base.ScoredIP{IP: sip.IP, Score: currentScore}, nil
		}
	}

	return nil, fmt.Errorf("没有找到满足阈值 %.2f 的 IP", threshold)
}
