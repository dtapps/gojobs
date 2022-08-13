package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"time"
)

// Lock 上锁
func (j *JobsGorm) Lock(ctx context.Context, info jobs_gorm_model.Task, id any) (string, error) {
	return j.lockClient.Lock(ctx, fmt.Sprintf("%s%s%v%s%v", j.config.lockKeyPrefix, j.config.lockKeySeparator, info.Type, j.config.lockKeySeparator, id), fmt.Sprintf("已在%s@%s机器上锁成功", j.config.insideIp, j.config.outsideIp), time.Duration(info.Frequency)*3*time.Second)
}

// Unlock Lock 解锁
func (j *JobsGorm) Unlock(ctx context.Context, info jobs_gorm_model.Task, id any) error {
	return j.lockClient.Unlock(ctx, fmt.Sprintf("%s%s%v%s%v", j.config.lockKeyPrefix, j.config.lockKeySeparator, info.Type, j.config.lockKeySeparator, id))
}

// LockForever 永远上锁
func (j *JobsGorm) LockForever(ctx context.Context, info jobs_gorm_model.Task, id any) (string, error) {
	return j.lockClient.LockForever(ctx, fmt.Sprintf("%s%s%v%s%v", j.config.lockKeyPrefix, j.config.lockKeySeparator, info.Type, j.config.lockKeySeparator, id), fmt.Sprintf("已在%s@%s机器永远上锁成功", j.config.insideIp, j.config.outsideIp))
}
