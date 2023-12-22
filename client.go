package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/golog"
)

// 前缀
// lockKeyPrefix 锁Key前缀 xxx_lock
// lockKeySeparator 锁Key分隔符 :
// cornKeyPrefix 任务Key前缀 xxx_cron
// cornKeyCustom 任务Key自定义 xxx_cron_自定义  xxx_cron_自定义_*
type redisPrefixFun func() (lockKeyPrefix, lockKeySeparator, cornKeyPrefix, cornKeyCustom string)

// ClientConfig 实例配置
type ClientConfig struct {
	GormClientFun  dorm.GormClientFun  // 数据库驱动
	RedisClientFun dorm.RedisClientFun // 数据库驱动
	RedisPrefixFun redisPrefixFun      // 前缀
	CurrentIp      string              // 当前IP
}

// Client 实例
type Client struct {
	gormClient *dorm.GormClient // 数据库
	config     struct {
		systemInsideIp  string // 内网ip
		systemOutsideIp string // 外网ip
	}
	cache struct {
		redisClient      *dorm.RedisClient     // 数据库
		redisLockClient  *dorm.RedisClientLock // 锁服务
		lockKeyPrefix    string                // 锁Key前缀 xxx_lock
		lockKeySeparator string                // 锁Key分隔符 :
		cornKeyPrefix    string                // 任务Key前缀 xxx_cron
		cornKeyCustom    string                // 任务Key自定义
	}
	slog struct {
		status bool        // 状态
		client *golog.SLog // 日志服务
	}
	runSlog struct {
		status bool        // 状态
		client *golog.SLog // 日志服务
	}
}

// NewClient 创建实例
func NewClient(config *ClientConfig) (*Client, error) {

	var ctx = context.Background()

	c := &Client{}

	if config.CurrentIp != "" && config.CurrentIp != "0.0.0.0" {
		c.config.systemOutsideIp = config.CurrentIp
	}
	c.config.systemOutsideIp = goip.IsIp(c.config.systemOutsideIp)
	if c.config.systemOutsideIp == "" {
		return nil, currentIpNoConfig
	}
	c.config.systemInsideIp = goip.GetInsideIp(ctx)

	// 配置缓存
	redisClient := config.RedisClientFun()
	if redisClient != nil && redisClient.GetDb() != nil {
		c.cache.redisClient = redisClient
		c.cache.redisLockClient = c.cache.redisClient.NewLock()
	} else {
		return nil, redisPrefixFunNoConfig
	}

	// 配置缓存前缀
	c.cache.lockKeyPrefix, c.cache.lockKeySeparator, c.cache.cornKeyPrefix, c.cache.cornKeyCustom = config.RedisPrefixFun()
	if c.cache.lockKeyPrefix == "" || c.cache.lockKeySeparator == "" || c.cache.cornKeyPrefix == "" || c.cache.cornKeyCustom == "" {
		return nil, redisPrefixFunNoConfig
	}

	// 配置关系数据库
	gormClient := config.GormClientFun()
	if gormClient != nil && gormClient.GetDb() != nil {
		c.gormClient = gormClient

		c.autoMigrateTask(ctx)
	} else {
		return nil, gormClientFunNoConfig
	}

	return c, nil
}
