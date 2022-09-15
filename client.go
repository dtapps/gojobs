package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/goip"
	"go.dtapp.net/golog"
	"log"
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
		systemHostName    string // 主机名
		systemInsideIp    string // 内网ip
		systemOs          string // 系统类型
		systemArch        string // 系统架构
		systemCpuQuantity int    // cpu核数
		goVersion         string // go版本
		sdkVersion        string // sdk版本
		systemMacAddrS    string // Mac地址
		systemOutsideIp   string // 外网ip
		debug             bool   // 日志开关
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

	// 配置外网ip
	if config.CurrentIp == "" {
		config.CurrentIp = goip.GetOutsideIp(ctx)
	}
	if config.CurrentIp != "" && config.CurrentIp != "0.0.0.0" {
		c.config.systemOutsideIp = config.CurrentIp
	}

	if c.config.debug {
		log.Printf("[gojobs]配置外网ip成功：%+v\n", c.config.systemOutsideIp)
	}

	// 配置缓存
	redisClient := config.RedisClientFun()
	if redisClient != nil && redisClient.Db != nil {
		c.cache.redisClient = redisClient
		c.cache.redisLockClient = c.cache.redisClient.NewLock()
	} else {
		return nil, redisPrefixFunNoConfig
	}

	if c.config.debug {
		log.Printf("[gojobs]配置缓存成功：%+v\n", c.cache)
	}

	// 配置缓存前缀
	c.cache.lockKeyPrefix, c.cache.lockKeySeparator, c.cache.cornKeyPrefix, c.cache.cornKeyCustom = config.RedisPrefixFun()
	if c.cache.lockKeyPrefix == "" || c.cache.lockKeySeparator == "" || c.cache.cornKeyPrefix == "" || c.cache.cornKeyCustom == "" {
		return nil, redisPrefixFunNoConfig
	}

	if c.config.debug {
		log.Printf("[gojobs]配置缓存前缀成功：%+v\n", c.cache)
	}

	// 配置信息
	c.setConfig(ctx)

	if c.config.debug {
		log.Printf("[gojobs]配置信息成功：%+v\n", c.config)
	}

	// 配置关系数据库
	gormClient := config.GormClientFun()
	if gormClient != nil && gormClient.Db != nil {
		c.db.gormClient = gormClient

		c.autoMigrateTask(ctx)
	} else {
		return nil, gormClientFunNoConfig
	}

	if c.config.debug {
		log.Printf("[gojobs]配置关系数据库成功：%+v\n", c.db)
	}

	// 配置非关系数据库
	mongoClient, databaseName := config.MongoClientFun()
	if mongoClient != nil && mongoClient.Db != nil {
		c.db.mongoClient = mongoClient
		if databaseName == "" {
			return nil, mongoClientFunNoConfig
		}
		c.db.mongoDatabaseName = databaseName

		c.mongoCreateCollectionTask(ctx)
		c.mongoCreateIndexesTask(ctx)
	} else {
		return nil, mongoClientFunNoConfig
	}

	if c.config.debug {
		log.Printf("[gojobs]配置非关系数据库成功：%+v\n", c.db)
	}

	if c.config.debug {
		log.Printf("[gojobs]创建实例成功：%+v\n", c)
	}

	return c, nil
}
