package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gorequest"
	"go.dtapp.net/gotime"
	"time"
)

// GetRedisKeyName 获取Redis键名
func GetRedisKeyName(taskType string) string {
	return "task:run:" + taskType
}

// SetRedisKeyValue 返回设置Redis键值
func SetRedisKeyValue(ctx context.Context, taskType string) (context.Context, string, any, time.Duration) {
	return ctx,
		GetRedisKeyName(taskType),
		fmt.Sprintf(
			"%s-%s-%s-%s",
			gotime.Current().SetFormat(gotime.DateTimeZhFormat),
			TraceGetTraceID(ctx),
			TraceGetSpanID(ctx),
			gorequest.GetRequestIDContext(ctx),
		),
		0
}

// SetRedisKeyValueExpiration 返回设置Redis键值，有过分时间
func SetRedisKeyValueExpiration(ctx context.Context, taskType string, expiration int64) (context.Context, string, any, time.Duration) {
	return ctx,
		GetRedisKeyName(taskType),
		fmt.Sprintf(
			"%s-%s-%s-%s",
			gotime.Current().SetFormat(gotime.DateTimeZhFormat),
			TraceGetTraceID(ctx),
			TraceGetSpanID(ctx),
			gorequest.GetRequestIDContext(ctx),
		),
		time.Duration(expiration)
}
