package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"time"
)

// LockId 上锁
func (j *JobsGorm) LockId(ctx context.Context, info jobs_gorm_model.Task) (string, error) {
	return j.lockClient.Lock(ctx, fmt.Sprintf("%s%s%v%s%v", j.config.lockKeyPrefix, j.config.lockKeySeparator, info.Type, j.config.lockKeySeparator, info.Id), fmt.Sprintf("已在%s@%s机器上锁成功", j.config.insideIp, j.config.outsideIp), time.Duration(info.Frequency)*3*time.Second)
}

// UnlockId Lock 解锁
func (j *JobsGorm) UnlockId(ctx context.Context, info jobs_gorm_model.Task) error {
	return j.lockClient.Unlock(ctx, fmt.Sprintf("%s%s%v%s%v", j.config.lockKeyPrefix, j.config.lockKeySeparator, info.Type, j.config.lockKeySeparator, info.Id))
}

// LockForeverId 永远上锁
func (j *JobsGorm) LockForeverId(ctx context.Context, info jobs_gorm_model.Task) (string, error) {
	return j.lockClient.LockForever(ctx, fmt.Sprintf("%s%s%v%s%v", j.config.lockKeyPrefix, j.config.lockKeySeparator, info.Type, j.config.lockKeySeparator, info.Id), fmt.Sprintf("已在%s@%s机器永远上锁成功", j.config.insideIp, j.config.outsideIp))
}
