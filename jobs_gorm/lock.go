package jobs_gorm

import (
	"fmt"
	"go.dtapp.net/goredis"
	"time"
)

// Lock 上锁
func (jobsGorm *JobsGorm) Lock(info Task, id any) string {
	cacheName := fmt.Sprintf("cron:%v:%v", info.Type, id)
	judgeCache := jobsGorm.Redis.NewStringOperation().Get(cacheName).UnwrapOr("")
	if judgeCache != "" {
		return judgeCache
	}
	jobsGorm.Redis.NewStringOperation().Set(cacheName, fmt.Sprintf("已在%v机器上锁成功", jobsGorm.OutsideIp), goredis.WithExpire(time.Millisecond*time.Duration(info.Frequency)*3))
	return ""
}

// Unlock Lock 解锁
func (jobsGorm *JobsGorm) Unlock(info Task, id any) {
	cacheName := fmt.Sprintf("cron:%v:%v", info.Type, id)
	jobsGorm.Redis.NewStringOperation().Del(cacheName)
}

// LockForever 永远上锁
func (jobsGorm *JobsGorm) LockForever(info Task, id any) string {
	cacheName := fmt.Sprintf("cron:%v:%v", info.Type, id)
	judgeCache := jobsGorm.Redis.NewStringOperation().Get(cacheName).UnwrapOr("")
	if judgeCache != "" {
		return judgeCache
	}
	jobsGorm.Redis.NewStringOperation().Set(cacheName, fmt.Sprintf("已在%v机器永远上锁成功", jobsGorm.OutsideIp))
	return ""
}
