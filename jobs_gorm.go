package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/golog"
	"runtime"
)

type JobsGormConfig struct {
	GormClient       *dorm.GormClient  // 数据库驱动
	RedisClient      *dorm.RedisClient // 缓存数据库驱动
	LogClient        *golog.ZapLog     // 日志驱动
	LogDebug         bool              // 日志开关
	CurrentIp        string            // 当前ip
	LockKeyPrefix    string            // 锁Key前缀 xxx_lock
	LockKeySeparator string            // 锁Key分隔符 :
	CornKeyPrefix    string            // 任务Key前缀 xxx_cron
	CornKeyCustom    string            // 任务Key自定义 xxx_cron_自定义  xxx_cron_自定义_*
}

// JobsGorm Gorm数据库驱动
type JobsGorm struct {
	gormClient  *dorm.GormClient      // 数据库服务
	redisClient *dorm.RedisClient     // 缓存服务
	lockClient  *dorm.RedisClientLock // 锁服务
	logClient   *golog.ZapLog         // 日志服务
	config      struct {
		logDebug         bool   // 日志开关
		runVersion       string // 运行版本
		os               string // 系统类型
		arch             string // 系统架构
		maxProCs         int    // CPU核数
		version          string // GO版本
		macAddrS         string // Mac地址
		insideIp         string // 内网ip
		outsideIp        string // 外网ip
		lockKeyPrefix    string // 锁Key前缀 xxx_lock
		lockKeySeparator string // 锁Key分隔符 :
		cornKeyPrefix    string // 任务Key前缀 xxx_cron
		cornKeyCustom    string // 任务Key自定义
	}
}

// NewJobsGorm 初始化
func NewJobsGorm(config *JobsGormConfig) (*JobsGorm, error) {

	// 判断
	if config.LockKeyPrefix == "" {
		return nil, errors.New("需要配置锁Key前缀")
	}
	if config.LockKeySeparator == "" {
		return nil, errors.New("需要配置锁Key分隔符")
	}
	if config.CornKeyPrefix == "" {
		return nil, errors.New("需要配置任务Key前缀")
	}
	if config.CornKeyCustom == "" {
		return nil, errors.New("需要配置任务Key自定义")
	}
	if config.CurrentIp == "" {
		return nil, errors.New("需要配置当前的IP")
	}

	if config.GormClient.Db == nil {
		return nil, errors.New("需要配置数据库驱动")
	}
	if config.RedisClient.Db == nil {
		return nil, errors.New("需要配置缓存数据库驱动")
	}

	c := &JobsGorm{}
	c.gormClient = config.GormClient
	c.redisClient = config.RedisClient
	c.lockClient = c.redisClient.NewLock()
	c.logClient = config.LogClient

	c.config.outsideIp = config.CurrentIp
	c.config.lockKeyPrefix = config.LockKeyPrefix
	c.config.lockKeySeparator = config.LockKeySeparator
	c.config.cornKeyPrefix = config.CornKeyPrefix
	c.config.cornKeyCustom = config.CornKeyCustom
	c.config.logDebug = config.LogDebug

	// 配置信息
	c.config.runVersion = Version
	c.config.os = runtime.GOOS
	c.config.arch = runtime.GOARCH
	c.config.maxProCs = runtime.GOMAXPROCS(0)
	c.config.version = runtime.Version()
	c.config.macAddrS = goarray.TurnString(goip.GetMacAddr(context.Background()))
	c.config.insideIp = goip.GetInsideIp(context.Background())

	// 创建模型
	err := c.gormClient.Db.AutoMigrate(
		&jobs_gorm_model.Task{},
		&jobs_gorm_model.TaskLog{},
		&jobs_gorm_model.TaskLogRun{},
		&jobs_gorm_model.TaskIp{},
	)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("创建任务模型失败：%v\n", err))
	}

	if c.config.logDebug == true {
		c.logClient.Infof(context.Background(), "[jobs.NewJobsGorm]%+v", c.config)
	}

	return c, nil
}
