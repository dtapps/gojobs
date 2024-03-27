package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/goip"
	"go.dtapp.net/gostring"
	"log/slog"
	"math/rand"
	"time"
)

// GetIssueAddress 获取下发地址
// workers 在线列表
// v 任务信息
// ---
// address 下发地址
// err 错误信息
func (c *Client) GetIssueAddress(ctx context.Context, workers []string, v *GormModelTask) (string, error) {
	var (
		currentIP       = ""    // 当前Ip
		appointIPStatus = false // 指定Ip状态
	)

	if v.SpecifyIP != "" {
		v.SpecifyIP = goip.IsIp(v.SpecifyIP)
	}

	// 赋值ip
	if v.SpecifyIP != "" && v.SpecifyIP != SpecifyIpNull {
		currentIP = v.SpecifyIP
		appointIPStatus = true
	}

	// 只有一个客户端在线
	if len(workers) == 1 {
		if appointIPStatus {
			// 判断是否指定某ip执行
			if gostring.Contains(workers[0], currentIP) {
				if c.slog.status {
					slog.InfoContext(ctx, fmt.Sprintf("只有一个客户端在线，指定某ip执行：%v %v", workers[0], currentIP))
				}
				return workers[0], nil
			}
			return "", errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIP))
		}
		return workers[0], nil
	}

	// 优先处理指定某ip执行
	if appointIPStatus {
		for wk, wv := range workers {
			if gostring.Contains(wv, currentIP) {
				if c.slog.status {
					slog.InfoContext(ctx, fmt.Sprintf("优先处理指定某ip执行：%v %v", workers[wk], currentIP))
				}
				return workers[wk], nil
			}
		}
		return "", errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIP))
	} else {
		// 随机返回一个
		address := workers[c.random(0, len(workers))]
		if address == "" {
			return address, errors.New("获取执行的客户端异常")
		}
		if c.slog.status {
			slog.InfoContext(ctx, fmt.Sprintf("随机返回一个：%v %v", address, currentIP))
		}
		return address, nil
	}
}

// GetSubscribeClientList 获取在线的客户端
func (c *Client) GetSubscribeClientList(ctx context.Context) (client []string, err error) {

	// 查询活跃的channel
	client, err = c.redisConfig.client.PubSubChannels(ctx, c.redisConfig.cornKeyPrefix+"_*").Result()
	if err != nil {
		if c.slog.status {
			slog.InfoContext(ctx, fmt.Sprintf("获取在线的客户端失败：%s，%v", c.redisConfig.cornKeyPrefix+"_*", err))
		}
	}

	return client, err
}

// 随机返回一个
// min最小
// max最大
func (c *Client) random(min, max int) int {
	if max-min <= 0 {
		return 0
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
