package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/golog"
)

// client *dorm.GormClient
type gormClientFun func() *dorm.GormClient

// client *dorm.MongoClient
// databaseName string
type mongoClientFun func() (*dorm.MongoClient, string)

// client *dorm.RedisClient
type redisClientFun func() *dorm.RedisClient

// 前缀
// lockKeyPrefix 锁Key前缀 xxx_lock
// lockKeySeparator 锁Key分隔符 :
// cornKeyPrefix 任务Key前缀 xxx_cron
// cornKeyCustom 任务Key自定义 xxx_cron_自定义  xxx_cron_自定义_*
type redisPrefixFun func() (lockKeyPrefix, lockKeySeparator, cornKeyPrefix, cornKeyCustom string)

// ClientConfig 实例配置
type ClientConfig struct {
	GormClientFun  gormClientFun  // 数据库驱动
	MongoClientFun mongoClientFun // 数据库驱动
	RedisClientFun redisClientFun // 数据库驱动
	RedisPrefixFun redisPrefixFun // 前缀
	Debug          bool           // 日志开关
	ZapLog         *golog.ZapLog  // 日志服务
	CurrentIp      string         // 当前ip
}

// Client 实例
type Client struct {
	zapLog *golog.ZapLog // 日志服务
	config struct {
		debug      bool   // 日志开关
		runVersion string // 运行版本
		os         string // 系统类型
		arch       string // 系统架构
		maxProCs   int    // CPU核数
		version    string // GO版本
		macAddrS   string // Mac地址
		insideIp   string // 内网ip
		outsideIp  string // 外网ip
	}
	cache struct {
		redisClient      *dorm.RedisClient     // 数据库
		redisLockClient  *dorm.RedisClientLock // 锁服务
		lockKeyPrefix    string                // 锁Key前缀 xxx_lock
		lockKeySeparator string                // 锁Key分隔符 :
		cornKeyPrefix    string                // 任务Key前缀 xxx_cron
		cornKeyCustom    string                // 任务Key自定义
	}
	db struct {
		gormClient        *dorm.GormClient  // 数据库
		mongoClient       *dorm.MongoClient // 数据库
		mongoDatabaseName string            // 数据库名
	}
}

// NewClient 创建实例
func NewClient(config *ClientConfig) (*Client, error) {

	var ctx = context.Background()

	c := &Client{}

	c.zapLog = config.ZapLog

	c.config.debug = config.Debug

	if config.CurrentIp == "" {
		return nil, currentIpNoConfig
	}
	c.config.outsideIp = config.CurrentIp

	// 配置信息
	c.setConfig(ctx)

	gormClient := config.GormClientFun()
	if gormClient != nil && gormClient.Db != nil {
		c.db.gormClient = gormClient

		c.autoMigrateTask()
		c.autoMigrateTaskIp()
		c.autoMigrateTaskLog()
		c.autoMigrateTaskLogRun()
	} else {
		return nil, gormClientFunNoConfig
	}

	mongoClient, databaseName := config.MongoClientFun()
	if mongoClient != nil && mongoClient.Db != nil {
		c.db.mongoClient = mongoClient
		if databaseName == "" {
			return nil, mongoClientFunNoConfig
		}
		c.db.mongoDatabaseName = databaseName

		c.mongoCreateCollectionTask(ctx)
		c.mongoCreateIndexesTask(ctx)
		c.mongoCreateIndexesTaskIp(ctx)
		c.mongoCreateCollectionTaskLog(ctx)
		c.mongoCreateIndexesTaskLog(ctx)
		c.mongoCreateCollectionTaskLogRun(ctx)
		c.mongoCreateIndexesTaskLogRun(ctx)
		c.mongoCreateCollectionTaskIssueRecord(ctx)
		c.mongoCreateCollectionTaskReceiveRecord(ctx)
	}

	redisClient := config.RedisClientFun()
	if redisClient != nil && redisClient.Db != nil {
		c.cache.redisClient = redisClient
		c.cache.redisLockClient = c.cache.redisClient.NewLock()
	}

	c.cache.lockKeyPrefix, c.cache.lockKeySeparator, c.cache.cornKeyPrefix, c.cache.cornKeyCustom = config.RedisPrefixFun()
	if c.cache.lockKeyPrefix == "" || c.cache.lockKeySeparator == "" || c.cache.cornKeyPrefix == "" || c.cache.cornKeyCustom == "" {
		return nil, redisPrefixFunNoConfig
	}

	return c, nil
}
