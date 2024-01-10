package gojobs

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.dtapp.net/golog"
	"gorm.io/gorm"
)

// 前缀
// lockKeyPrefix 锁Key前缀 xxx_lock
// lockKeySeparator 锁Key分隔符 :
// cornKeyPrefix 任务Key前缀 xxx_cron
// cornKeyCustom 任务Key自定义 xxx_cron_自定义  xxx_cron_自定义_*
type redisPrefixFun func() (lockKeyPrefix, lockKeySeparator, cornKeyPrefix, cornKeyCustom string)

// ClientConfig 实例配置
type ClientConfig struct {
	GormClient     *gorm.DB       // 数据库驱动
	RedisClient    *redis.Client  // 数据库驱动
	RedisPrefixFun redisPrefixFun // 前缀
	CurrentIp      string         // 当前IP
}

// Client 实例
type Client struct {
	gormClient *gorm.DB // 数据库
	config     struct {
		systemHostname      string  // 主机名
		systemOs            string  // 系统类型
		systemVersion       string  // 系统版本
		systemKernel        string  // 系统内核
		systemKernelVersion string  // 系统内核版本
		systemUpTime        uint64  // 系统运行时间
		systemBootTime      uint64  // 系统开机时间
		cpuCores            int     // CPU核数
		cpuModelName        string  // CPU型号名称
		cpuMhz              float64 // CPU兆赫
		systemInsideIP      string  // 内网IP
		systemOutsideIP     string  // 外网IP
		goVersion           string  // go版本
		sdkVersion          string  // sdk版本
		logVersion          string  // log版本
		redisSdkVersion     string  // redis版本
	}
	cache struct {
		redisClient      *redis.Client // 数据库
		lockKeyPrefix    string        // 锁Key前缀 xxx_lock
		lockKeySeparator string        // 锁Key分隔符 :
		cornKeyPrefix    string        // 任务Key前缀 xxx_cron
		cornKeyCustom    string        // 任务Key自定义
	}
	slog struct {
		status bool        // 状态
		client *golog.SLog // 日志服务
	}
}

// NewClient 创建实例
func NewClient(config *ClientConfig) (*Client, error) {

	var ctx = context.Background()

	c := &Client{}

	if config.CurrentIp == "" || config.CurrentIp == "0.0.0.0" {
		return nil, currentIpNoConfig
	}

	// 配置缓存
	redisClient := config.RedisClient
	if redisClient != nil {
		c.cache.redisClient = redisClient
	} else {
		return nil, redisPrefixFunNoConfig
	}

	// 配置缓存前缀
	c.cache.lockKeyPrefix, c.cache.lockKeySeparator, c.cache.cornKeyPrefix, c.cache.cornKeyCustom = config.RedisPrefixFun()
	if c.cache.lockKeyPrefix == "" || c.cache.lockKeySeparator == "" || c.cache.cornKeyPrefix == "" || c.cache.cornKeyCustom == "" {
		return nil, redisPrefixFunNoConfig
	}

	// 配置信息
	c.setConfig(ctx, config.CurrentIp)

	// 配置关系数据库
	gormClient := config.GormClient
	if gormClient != nil {
		c.gormClient = gormClient

		c.autoMigrateTask(ctx)
	} else {
		return nil, gormClientFunNoConfig
	}

	return c, nil
}
