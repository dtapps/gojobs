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
	MongoClientFun dorm.MongoClientFun // 数据库驱动
	RedisClientFun dorm.RedisClientFun // 数据库驱动
	RedisPrefixFun redisPrefixFun      // 前缀
	ZapLog         *golog.ZapLog       // 日志服务
	CurrentIp      string              // 当前ip
}

// Client 实例
type Client struct {
	gormClient  *dorm.GormClient  // 数据库
	mongoClient *dorm.MongoClient // 数据库
	zapLog      *golog.ZapLog     // 日志服务
	config      struct {
		systemHostname      string  // 主机名
		systemOs            string  // 系统类型
		systemVersion       string  // 系统版本
		systemKernel        string  // 系统内核
		systemKernelVersion string  // 系统内核版本
		systemBootTime      uint64  // 系统开机时间
		cpuCores            int     // CPU核数
		cpuModelName        string  // CPU型号名称
		cpuMhz              float64 // CPU兆赫
		systemInsideIp      string  // 内网ip
		systemOutsideIp     string  // 外网ip
		goVersion           string  // go版本
		sdkVersion          string  // sdk版本
	}
	cache struct {
		redisClient      *dorm.RedisClient     // 数据库
		redisLockClient  *dorm.RedisClientLock // 锁服务
		lockKeyPrefix    string                // 锁Key前缀 xxx_lock
		lockKeySeparator string                // 锁Key分隔符 :
		cornKeyPrefix    string                // 任务Key前缀 xxx_cron
		cornKeyCustom    string                // 任务Key自定义
	}
	mongoConfig struct {
		stats        bool   // 状态
		databaseName string // 库名
	}
}

// NewClient 创建实例
func NewClient(config *ClientConfig) (*Client, error) {

	var ctx = context.Background()

	c := &Client{}

	c.zapLog = config.ZapLog

	if config.CurrentIp != "" && config.CurrentIp != "0.0.0.0" {
		c.config.systemOutsideIp = config.CurrentIp
	}
	c.config.systemOutsideIp = goip.IsIp(c.config.systemOutsideIp)
	if c.config.systemOutsideIp == "" {
		return nil, currentIpNoConfig
	}

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

	// 配置信息
	c.setConfig(ctx)

	// 配置关系数据库
	gormClient := config.GormClientFun()
	if gormClient != nil && gormClient.GetDb() != nil {
		c.gormClient = gormClient

		c.autoMigrateTask(ctx)
		c.autoMigrateTaskLog(ctx)
	} else {
		return nil, gormClientFunNoConfig
	}

	// 配置非关系数据库
	mongoClient, databaseName := config.MongoClientFun()
	if mongoClient != nil && mongoClient.GetDb() != nil {
		c.mongoClient = mongoClient

		if databaseName == "" {
			return nil, mongoClientFunNoConfig
		} else {
			c.mongoConfig.databaseName = databaseName
		}

		TaskLog{}.createCollection(ctx, c.zapLog, c.mongoClient, c.mongoConfig.databaseName)
		TaskLog{}.createIndexes(ctx, c.zapLog, c.mongoClient, c.mongoConfig.databaseName)

		c.mongoConfig.stats = true
	}

	return c, nil
}
