package gojobs

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
)

// Publish 发布
// ctx 上下文
// channel 频道
// message 消息
func (c *Client) Publish(ctx context.Context, channel string, message interface{}) error {
	publish, err := c.redisConfig.client.Publish(ctx, channel, message).Result()
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("发布失败：%s %s %v %s", channel, message, publish, err))
		}
	}
	return err
}

type SubscribeResult struct {
	err     error
	Message *redis.PubSub
}

// Subscribe 订阅
func (c *Client) Subscribe(ctx context.Context) SubscribeResult {
	return SubscribeResult{
		Message: c.redisConfig.client.Subscribe(ctx, c.redisConfig.cornKeyPrefix+"_"+c.redisConfig.cornKeyCustom),
	}
}

// PSubscribe 订阅，支持通配符匹配(ch_user_*)
func (c *Client) PSubscribe(ctx context.Context) SubscribeResult {
	return SubscribeResult{
		Message: c.redisConfig.client.PSubscribe(ctx, c.redisConfig.cornKeyPrefix+"_"+c.redisConfig.cornKeyCustom+"_*"),
	}
}
