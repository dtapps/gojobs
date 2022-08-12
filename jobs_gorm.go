package gojobs

import (
	"errors"
	"fmt"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/golock"
	"log"
	"runtime"
)

type JobsGormConfig struct {
	GormClient    *dorm.GormClient  // 数据库驱动
	RedisClient   *dorm.RedisClient // 缓存数据库驱动
	CurrentIp     string            // 当前ip
	LockPrefix    string            // 锁Key前缀 xxx_lock
	LockSeparator string            // 锁分隔符 xxx_lock
	CornPrefix    string            // 任务前缀 xxx_cron
	Debug         bool              // 调试
}

// JobsGorm Gorm数据库驱动
type JobsGorm struct {
	gormClient  *dorm.GormClient  // 数据库驱动
	redisClient *dorm.RedisClient // 缓存驱动
	lockClient  *golock.LockRedis // 锁驱动
	config      struct {
		debug           bool   // 调试
		runVersion      string // 运行版本
		os              string // 系统类型
		arch            string // 系统架构
		maxProCs        int    // CPU核数
		version         string // GO版本
		macAddrS        string // Mac地址
		insideIp        string // 内网ip
		outsideIp       string // 外网ip
		lockPrefix      string // 锁Key前缀
		lockSeparator   string // 锁分隔符
		cornPrefix      string // 任务key前缀
		cornKeyIp       string // 任务key
		cornKeyChannel  string // 任务频道key(任务key+ip)
		cornKeyChannels string // 任务频道key通配符匹配(任务key+ip+_*)
	}
}

// NewJobsGorm 初始化
func NewJobsGorm(config *JobsGormConfig) (*JobsGorm, error) {

	// 判断
	if config.LockPrefix == "" {
		return nil, errors.New("需要配置锁Key前缀")
	}
	if config.LockSeparator == "" {
		return nil, errors.New("需要配置锁分隔符")
	}
	if config.CornPrefix == "" {
		return nil, errors.New("需要配置任务前缀")
	}
	if config.CurrentIp == "" {
		return nil, errors.New("需要配置当前的IP")
	}
	if config.GormClient == nil {
		return nil, errors.New("需要配置数据库驱动")
	}
	if config.RedisClient == nil {
		return nil, errors.New("需要配置缓存数据库驱动")
	}

	c := &JobsGorm{}
	c.gormClient = config.GormClient
	c.redisClient = config.RedisClient
	c.config.outsideIp = config.CurrentIp
	c.config.lockPrefix = config.LockPrefix
	c.config.lockSeparator = config.LockSeparator
	c.config.cornPrefix = config.CornPrefix
	c.config.debug = config.Debug

	// 锁
	c.lockClient = golock.NewLockRedis(c.redisClient)

	// 配置信息
	c.config.runVersion = Version
	c.config.os = runtime.GOOS
	c.config.arch = runtime.GOARCH
	c.config.maxProCs = runtime.GOMAXPROCS(0)
	c.config.version = runtime.Version()
	c.config.macAddrS = goarray.TurnString(goip.GetMacAddr())
	c.config.insideIp = goip.GetInsideIp()

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

	c.config.cornKeyIp = c.getCornKeyIp()
	c.config.cornKeyChannel = c.getCornKeyChannel()
	c.config.cornKeyChannels = c.getCornKeyChannels()

	if c.config.cornKeyIp == "" {
		return nil, errors.New(fmt.Sprintf("没有配置 cornKeyIp：%s", c.config.cornKeyIp))
	}
	if c.config.cornKeyChannel == "" {
		return nil, errors.New(fmt.Sprintf("没有配置 cornKeyChannel：%s", c.config.cornKeyChannel))
	}
	if c.config.cornKeyChannels == "" {
		return nil, errors.New(fmt.Sprintf("没有配置 cornKeyChannels：%s", c.config.cornKeyChannels))
	}

	if c.config.debug == true {
		log.Printf("JOBS配置：%+v\n", c.config)
	}

	return c, nil
}
