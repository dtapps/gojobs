package gojobs

import (
	"context"
	"github.com/go-redis/redis/v9"
)

// Publish 发布
// ctx 上下文
// channel 频道
// message 消息
func (j *JobsGorm) Publish(ctx context.Context, channel string, message interface{}) error {
	publish, err := j.redisClient.Publish(ctx, channel, message).Result()
	if j.config.logDebug == true {
		j.logClient.Infof(ctx, "[jobs.Publish] %s %s %v %s", channel, message, publish, err)
	}
	return err
}

type SubscribeResult struct {
	err     error
	Message *redis.PubSub
}

// Subscribe 订阅
func (j *JobsGorm) Subscribe(ctx context.Context) SubscribeResult {
	return SubscribeResult{
		Message: j.redisClient.Subscribe(ctx, j.config.cornKeyPrefix+"_"+j.config.cornKeyCustom),
	}
}

// PSubscribe 订阅，支持通配符匹配(ch_user_*)
func (j *JobsGorm) PSubscribe(ctx context.Context) SubscribeResult {
	return SubscribeResult{
		Message: j.redisClient.PSubscribe(ctx, j.config.cornKeyPrefix+"_"+j.config.cornKeyCustom+"_*"),
	}
}

// PubSubChannels 查询活跃的channel
func (j *JobsGorm) PubSubChannels(ctx context.Context) []string {
	result, _ := j.redisClient.PubSubChannels(ctx, j.config.cornKeyPrefix+"_"+j.config.cornKeyCustom+"_*").Result()
	return result
}

// PubSubNumSub 查询指定的channel有多少个订阅者
func (j *JobsGorm) PubSubNumSub(ctx context.Context) map[string]int64 {
	result, _ := j.redisClient.PubSubNumSub(ctx, j.config.cornKeyPrefix+"_"+j.config.cornKeyCustom+"_*").Result()
	return result
}
