package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gostring"
	"go.dtapp.net/gotrace_id"
	"log"
	"math/rand"
	"time"
)

// GetIssueAddress 获取下发地址
// workers 在线列表
// v 任务信息
// ---
// address 下发地址
// err 错误信息
func (j *JobsGorm) GetIssueAddress(ctx context.Context, workers []string, v *jobs_gorm_model.Task) (string, error) {
	var (
		currentIp       = ""    // 当前Ip
		appointIpStatus = false // 指定Ip状态
		traceId         = gotrace_id.GetTraceIdContext(ctx)
	)

	// 赋值ip
	if v.SpecifyIp != "" {
		currentIp = v.SpecifyIp
		appointIpStatus = true
	}

	// 只有一个客户端在线
	if len(workers) == 1 {
		if appointIpStatus == true {
			// 判断是否指定某ip执行
			if gostring.Contains(workers[0], currentIp) == true {
				log.Println("[jobs.GetIssueAddress]只有一个客户端在线，指定某ip执行", traceId, workers[0], currentIp)
				return workers[0], nil
			}
			return "", errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIp))
		}
		return workers[0], nil
	}

	// 优先处理指定某ip执行
	if appointIpStatus == true {
		for wk, wv := range workers {
			if gostring.Contains(wv, currentIp) == true {
				log.Println("[jobs.GetIssueAddress]优先处理指定某ip执行", traceId, workers[wk], currentIp)
				return workers[wk], nil
			}
		}
		return "", errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIp))
	} else {
		// 随机返回一个
		address := workers[j.random(0, len(workers))]
		if address == "" {
			return address, errors.New("获取执行的客户端异常")
		}
		log.Println("[jobs.GetIssueAddress]随机返回一个", traceId, address, currentIp)
		return address, nil
	}
}

// GetSubscribeClientList 获取在线的客户端
func (j *JobsGorm) GetSubscribeClientList(ctx context.Context) ([]string, error) {

	if j.config.logDebug == true {
		j.logClient.Infof(ctx, "[jobs.GetSubscribeClientList] %s", j.config.cornKeyPrefix+"_*")
	}

	// 扫描
	client := j.redisClient.Keys(ctx, j.config.cornKeyPrefix+"_*")

	return client, nil
}

// 随机返回一个
// min最小
// max最大
func (j *JobsGorm) random(min, max int) int {
	if max-min <= 0 {
		return 0
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
