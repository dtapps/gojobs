package gojobs

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.dtapp.net/golog"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

// Client 实例
type Client struct {
	config struct {
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
		client           *gorm.DB // 数据库
		taskTableName    string   // 任务表名
		taskLogStatus    bool     // 任务日志状态
		taskLogTableName string   // 任务日志表名
	}
	mongoConfig struct {
		client                *mongo.Client // 数据库
		databaseName          string        // 库名
		taskLogStatus         bool          // 任务日志状态
		taskLogCollectionName string        // 任务日志集合名
	}
	slog struct {
		status bool        // 状态
		client *golog.SLog // 日志服务
	}
}

// NewClient 创建实例
func NewClient(ctx context.Context, currentIP string) (*Client, error) {

	c := &Client{}

	if currentIP == "" || currentIP == "0.0.0.0" {
		return nil, errors.New("请配置 CurrentIp")
	}

	// 配置信息
	c.setConfig(ctx, currentIP)

	return c, nil
}
