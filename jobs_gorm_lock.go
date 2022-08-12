package gojobs

import (
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
)

// Lock 上锁
func (j *JobsGorm) Lock(info jobs_gorm_model.Task, id any) (string, error) {
	return j.lockClient.Lock(fmt.Sprintf("%s%s%v%s%v", j.config.lockPrefix, j.config.lockSeparator, info.Type, j.config.lockSeparator, id), fmt.Sprintf("已在%s@%s机器上锁成功", j.config.insideIp, j.config.outsideIp), info.Frequency*3)
}

// Unlock Lock 解锁
func (j *JobsGorm) Unlock(info jobs_gorm_model.Task, id any) error {
	return j.lockClient.Unlock(fmt.Sprintf("%s%s%v%s%v", j.config.lockPrefix, j.config.lockSeparator, info.Type, j.config.lockSeparator, id))
}

// LockForever 永远上锁
func (j *JobsGorm) LockForever(info jobs_gorm_model.Task, id any) (string, error) {
	return j.lockClient.LockForever(fmt.Sprintf("%s%s%v%s%v", j.config.lockPrefix, j.config.lockSeparator, info.Type, j.config.lockSeparator, id), fmt.Sprintf("已在%s@%s机器永远上锁成功", j.config.insideIp, j.config.outsideIp))
}
