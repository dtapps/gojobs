package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/gotime"
	"log"
	"time"
)

func (c *Client) StartHandle(ctx context.Context, key any, overdue int64) error {
	status, err := c.redisConfig.client.Get(ctx, fmt.Sprintf("%v", key)).Result()
	if status != "" {
		return errors.New(fmt.Sprintf("【%v】上一次还在执行中", fmt.Sprintf("%v", key)))
	}

	err = c.redisConfig.client.Set(ctx, fmt.Sprintf("%v", key), gotime.Current().Format(), time.Duration(overdue)*time.Second).Err()
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("【%v】设置 %s", fmt.Sprintf("%v", key), err))
		}
	}

	return nil
}
func (c *Client) EndHandle(ctx context.Context, key any) {
	err := c.redisConfig.client.Del(ctx, fmt.Sprintf("%v", key)).Err()
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("【%v】删除 %s", fmt.Sprintf("%v", key), err))
		}
	}
}
