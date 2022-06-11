package jobs_gorm

import (
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs"
	"go.dtapp.net/goredis"
	"gorm.io/gorm"
	"runtime"
)

type JobsGorm struct {
	runVersion  string      // 运行版本
	os          string      // 系统类型
	arch        string      // 系统架构
	maxProCs    int         // CPU核数
	version     string      // GO版本
	macAddrS    string      // Mac地址
	insideIp    string      // 内网ip
	outsideIp   string      // 外网ip
	mainService int         // 主要服务
	Db          *gorm.DB    // 数据库
	Redis       goredis.App // 缓存数据库服务
}

func NewGorm(jobsGorm JobsGorm, mainService int) *JobsGorm {
	jobsGorm.runVersion = gojobs.Version
	jobsGorm.os = runtime.GOOS
	jobsGorm.arch = runtime.GOARCH
	jobsGorm.maxProCs = runtime.GOMAXPROCS(0)
	jobsGorm.version = runtime.Version()
	jobsGorm.macAddrS = goarray.TurnString(goip.GetMacAddr())
	jobsGorm.insideIp = goip.GetInsideIp()
	jobsGorm.outsideIp = goip.GetOutsideIp()
	jobsGorm.mainService = mainService
	return &jobsGorm
}
