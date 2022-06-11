package gojobs

import (
	"errors"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm"
	"go.dtapp.net/goredis"
	"gorm.io/gorm"
)

type ConfigJobsGorm struct {
	MainService int             // 主要服务
	Db          *gorm.DB        // 数据库
	Redis       *goredis.Client // 缓存数据库服务
}

func NewJobsGorm(config *ConfigJobsGorm) *jobs_gorm.JobsGorm {

	var (
		jobsGorm = &jobs_gorm.JobsGorm{}
	)

	jobsGorm = jobs_gorm.NewGorm(jobs_gorm.JobsGorm{
		Db:    config.Db,
		Redis: config.Redis,
	}, config.MainService, Version)

	err := jobsGorm.Db.AutoMigrate(
		&jobs_gorm.Task{},
		&jobs_gorm.TaskLog{},
		&jobs_gorm.TaskLogRun{},
		&jobs_gorm.TaskIp{},
	)
	if err != nil {
		panic(errors.New(fmt.Sprintf("创建任务模型失败：%v\n", err)))
	}

	return jobsGorm
}
