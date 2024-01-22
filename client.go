package gojobs

import (
	"context"
	"github.com/redis/go-redis/v9"
	"go.dtapp.net/golog"
	"go.mongodb.org/mongo-driver/mongo"
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
	GormClient          *gorm.DB       // 关系数据库驱动
	GormTableName       string         // 关系数据库表名
	MongoClient         *mongo.Client  // MONGO数据库驱动
	MongoDatabaseName   string         // MONGO数据库名
	MongoCollectionName string         // MONGO数据库集合名
	RedisClient         *redis.Client  // Redis数据库驱动
	RedisPrefixFun      redisPrefixFun // Redis数据前缀
	CurrentIP           string         // 当前IP
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
		redisSdkVersion     string  // redisSdk版本
	}
	redisConfig struct {
		client           *redis.Client // 数据库
		lockKeyPrefix    string        // 锁Key前缀 xxx_lock
		lockKeySeparator string        // 锁Key分隔符 :
		cornKeyPrefix    string        // 任务Key前缀 xxx_cron
		cornKeyCustom    string        // 任务Key自定义
	}
	gormConfig struct {
		client *gorm.DB // 数据库
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

	if config.CurrentIP == "" || config.CurrentIP == "0.0.0.0" {
		return nil, currentIpNoConfig
	}

	// 配置缓存
	redisClient := config.RedisClient
	if redisClient != nil {
		c.redisConfig.client = redisClient
	} else {
		return nil, redisPrefixFunNoConfig
	}

	// 配置缓存前缀
	c.redisConfig.lockKeyPrefix, c.redisConfig.lockKeySeparator, c.redisConfig.cornKeyPrefix, c.redisConfig.cornKeyCustom = config.RedisPrefixFun()
	if c.redisConfig.lockKeyPrefix == "" || c.redisConfig.lockKeySeparator == "" || c.redisConfig.cornKeyPrefix == "" || c.redisConfig.cornKeyCustom == "" {
		return nil, redisPrefixFunNoConfig
	}

	// 配置信息
	c.setConfig(ctx, config.CurrentIP)

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
