package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gostring"
	"math/rand"
	"time"
)

// GetIssueAddress 获取下发地址
// workers 在线列表
// v 任务信息
// ---
// address 下发地址
// err 错误信息
func (c *Client) GetIssueAddress(ctx context.Context, workers []string, v *jobs_gorm_model.Task) (string, error) {
	var (
		currentIp       = ""    // 当前Ip
		appointIpStatus = false // 指定Ip状态
	)

	// 赋值ip
	if v.SpecifyIp != "" && v.SpecifyIp != SpecifyIpNull {
		currentIp = v.SpecifyIp
		appointIpStatus = true
	}

	// 只有一个客户端在线
	if len(workers) == 1 {
		if appointIpStatus {
			// 判断是否指定某ip执行
			if gostring.Contains(workers[0], currentIp) {
				c.zapLog.WithTraceId(ctx).Sugar().Info("[jobs]只有一个客户端在线，指定某ip执行：", workers[0], currentIp)
				return workers[0], nil
			}
			return "", errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIp))
		}
		return workers[0], nil
	}

	// 优先处理指定某ip执行
	if appointIpStatus {
		for wk, wv := range workers {
			if gostring.Contains(wv, currentIp) {
				c.zapLog.WithTraceId(ctx).Sugar().Info("[jobs]优先处理指定某ip执行：", workers[wk], currentIp)
				return workers[wk], nil
			}
		}
		return "", errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIp))
	} else {
		// 随机返回一个
		address := workers[c.random(0, len(workers))]
		if address == "" {
			return address, errors.New("获取执行的客户端异常")
		}
		c.zapLog.WithTraceId(ctx).Sugar().Info("[jobs]随机返回一个：", address, currentIp)
		return address, nil
	}
}

// GetSubscribeClientList 获取在线的客户端
func (c *Client) GetSubscribeClientList(ctx context.Context) (client []string, err error) {

	// 查询活跃的channel
	client, err = c.cache.redisClient.PubSubChannels(ctx, c.cache.cornKeyPrefix+"_*").Result()
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[jobs]获取在线的客户端失败：%s，%v", c.cache.cornKeyPrefix+"_*", err)
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
