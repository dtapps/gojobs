package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
	"time"
)

// LockId 上锁
func (c *Client) LockId(ctx context.Context, info jobs_gorm_model.Task) (string, error) {
	return c.lock(ctx, fmt.Sprintf("%s%s%v%s%v", c.cache.lockKeyPrefix, c.cache.lockKeySeparator, info.Type, c.cache.lockKeySeparator, info.ID), fmt.Sprintf("[LockId] 已在%s@%s机器上锁成功，时间：%v", c.config.systemInsideIP, c.config.systemOutsideIP, gotime.Current().Format()), time.Duration(info.Frequency)*3*time.Second)
}

// UnlockId 解锁
func (c *Client) UnlockId(ctx context.Context, info jobs_gorm_model.Task) error {
	return c.unlock(ctx, fmt.Sprintf("%s%s%v%s%v", c.cache.lockKeyPrefix, c.cache.lockKeySeparator, info.Type, c.cache.lockKeySeparator, info.ID))
}

// LockForeverId 永远上锁
func (c *Client) LockForeverId(ctx context.Context, info jobs_gorm_model.Task) (string, error) {
	return c.lockForever(ctx, fmt.Sprintf("%s%s%v%s%v", c.cache.lockKeyPrefix, c.cache.lockKeySeparator, info.Type, c.cache.lockKeySeparator, info.ID), fmt.Sprintf("[LockForeverId] 已在%s@%s机器永远上锁成功，时间：%v", c.config.systemInsideIP, c.config.systemOutsideIP, gotime.Current().Format()))
}
