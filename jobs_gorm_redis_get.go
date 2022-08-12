package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"log"
	"math/rand"
	"time"
)

// GetIssueAddress 获取下发地址
func (j *JobsGorm) GetIssueAddress(ctx context.Context, v *jobs_gorm_model.Task) (address string, err error) {
	var (
		currentIp       = ""    // 当前Ip
		appointIpStatus = false // 指定Ip状态
	)

	// 赋值ip
	if v.SpecifyIp != "" {
		currentIp = v.SpecifyIp
		appointIpStatus = true
	}

	workers, err := j.GetSubscribeClientList(ctx)
	if err != nil {
		return address, errors.New(fmt.Sprintf("获取在线客户端列表失败：%s", err.Error()))
	}
	if len(workers) <= 0 {
		return address, errors.New("没有客户端在线")
	}

	// 只有一个客户端在线
	if len(workers) == 1 {
		if appointIpStatus == true {
			// 判断是否指定某ip执行
			if currentIp == workers[0] {
				return j.config.cornPrefix + "_" + v.SpecifyIp, nil
			}
			return address, errors.New("执行的客户端不在线")
		}
		return j.config.cornPrefix + "_" + workers[0], nil
	}

	// 优先处理指定某ip执行
	if appointIpStatus == true {
		for _, wv := range workers {
			if currentIp == wv {
				return j.config.cornPrefix + "_" + wv, nil
			}
		}
		return address, errors.New("执行的客户端不在线")
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
func (j *JobsGorm) random(min, max int) int {
	if max-min <= 0 {
		return 0
	}
	rand.Seed(time.Now().Unix())
	return rand.Intn(max-min) + min
}
