package gojobs

import (
	"go.dtapp.net/gojobs/jobs_gorm"
	"go.dtapp.net/goredis"
	"gorm.io/gorm"
)

type ConfigJobsGorm struct {
	MainService int         // 主要服务
	Db          *gorm.DB    // 数据库
	Redis       goredis.App // 缓存数据库服务
}

func NewJobsGorm(config *ConfigJobsGorm) *jobs_gorm.JobsGorm {
	var (
		jobsGorm = &jobs_gorm.JobsGorm{}
	)
	jobsGorm = jobs_gorm.NewGorm(jobs_gorm.JobsGorm{
		Db:    config.Db,
		Redis: config.Redis,
	}, config.MainService, Version)
	return jobsGorm
}
