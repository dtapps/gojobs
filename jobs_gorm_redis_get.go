package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gostring"
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
func (j *JobsGorm) GetIssueAddress(workers []string, v *jobs_gorm_model.Task) (address string, err error) {
	var (
		currentIp       = ""    // 当前Ip
		appointIpStatus = false // 指定Ip状态
	)

	// 赋值ip
	if v.SpecifyIp != "" {
		currentIp = v.SpecifyIp
		appointIpStatus = true
	}

	//workers, err := j.GetSubscribeClientList(ctx)
	//if err != nil {
	//	return address, errors.New(fmt.Sprintf("获取在线客户端列表失败：%s", err.Error()))
	//}
	//if len(workers) <= 0 {
	//	return address, errors.New("没有客户端在线")
	//}

	// 只有一个客户端在线
	if len(workers) == 1 {
		if appointIpStatus == true {
			// 判断是否指定某ip执行
			if gostring.Contains(currentIp, workers[0]) == true {
				return j.config.cornPrefix + "_" + v.SpecifyIp, nil
			}
			return address, errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIp))
		}
		return j.config.cornPrefix + "_" + workers[0], nil
	}

	// 优先处理指定某ip执行
	if appointIpStatus == true {
		for wk, wv := range workers {
			if gostring.Contains(currentIp, wv) == true {
				return j.config.cornPrefix + "_" + workers[wk], nil
			}
		}
		return address, errors.New(fmt.Sprintf("需要执行的[%s]客户端不在线", currentIp))
	} else {
		// 随机返回一个
		zxIp := workers[j.random(0, len(workers))]
		if zxIp == "" {
			return address, errors.New("获取执行的客户端异常")
		}
		address = j.config.cornPrefix + "_" + zxIp
		return address, err
	}
}

// GetSubscribeClientList 获取在线的客户端
func (j *JobsGorm) GetSubscribeClientList(ctx context.Context) ([]string, error) {

	if j.config.debug == true {
		log.Printf("获取在线的客户端：%s\n", j.config.cornPrefix+"_*")
	}

	// 扫描
	values, err := j.redisClient.Keys(ctx, j.config.cornPrefix+"_*").Result()
	if err != nil {
		if err == errors.New("ERR wrong number of arguments for 'mget' command") {
			return []string{}, nil
		}
		return nil, errors.New(fmt.Sprintf("获取失败：%s", err.Error()))
	}

	client := make([]string, 0, len(values))
	for _, val := range values {
		client = append(client, val.(string))
	}

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
