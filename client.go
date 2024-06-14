package gojobs

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

// Client 实例
type Client struct {
	config struct {
		systemInsideIP  string // 内网IP
		systemOutsideIP string // 外网IP
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
	slog struct {
		status bool // 状态
	}
	trace bool       // OpenTelemetry链路追踪
	span  trace.Span // OpenTelemetry链路追踪
}

// NewClient 创建实例
func NewClient(ctx context.Context, currentIP string) (*Client, error) {
	c := &Client{}

	if currentIP == "" || currentIP == "0.0.0.0" {
		return nil, errors.New("请配置 CurrentIp")
	}

	// 配置信息
	c.setConfig(ctx, currentIP)

	c.trace = true
	return c, nil
}
