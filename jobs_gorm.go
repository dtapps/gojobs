package gojobs

import (
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_gorm"
	"go.dtapp.net/goredis"
	"gorm.io/gorm"
	"runtime"
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
	jobsGorm.RunVersion = Version
	jobsGorm.Os = runtime.GOOS
	jobsGorm.Arch = runtime.GOARCH
	jobsGorm.MaxProCs = runtime.GOMAXPROCS(0)
	jobsGorm.Version = runtime.Version()
	jobsGorm.MacAddrS = goarray.TurnString(goip.GetMacAddr())
	jobsGorm.InsideIp = goip.GetInsideIp()
	jobsGorm.OutsideIp = goip.GetOutsideIp()
	jobsGorm.MainService = config.MainService
	jobsGorm.Db = config.Db
	jobsGorm.Redis = config.Redis
	return jobsGorm
}
