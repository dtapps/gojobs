package gojobs

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type PubSubClient struct {
	client *redis.Client
}

func NewPubSub(ctx context.Context, client *redis.Client) *PubSubClient {
	return &PubSubClient{
		client: client,
	}
}

// Publish 发布
func (c *PubSubClient) Publish(ctx context.Context, channel string, message interface{}) error {
	return c.client.Publish(ctx, channel, message).Err()
}

// Subscribe 订阅
func (c *PubSubClient) Subscribe(ctx context.Context, channel ...string) *redis.PubSub {
	return c.client.Subscribe(ctx, channel...)
}

// PSubscribe 订阅，支持通配符匹配(ch_user_*)
func (c *PubSubClient) PSubscribe(ctx context.Context, channel ...string) *redis.PubSub {
	return c.client.PSubscribe(ctx, channel...)
}

// CustomRunSingleTask 运行单个任务
// ctx 链路追踪的上下文
// message 任务json编码
// executionCallback 执行任务回调函数 返回 runCode=状态 runDesc=描述
// updateCallback 执行更新回调函数
func (c *PubSubClient) CustomRunSingleTask(ctx context.Context, message string, executionCallback func(ctx context.Context, task *TaskCustomHelperTaskList) (err error)) {

	// 解析任务
	var task TaskCustomHelperTaskList
	err := gojson.Unmarshal([]byte(message), &task)
	if err != nil {
		return
	}

	// 启动OpenTelemetry链路追踪
	ctx, span := NewTraceStartSpan(ctx, "CustomRunSingleTask "+task.TaskName+" "+task.CustomID)

	span.SetAttributes(attribute.String("task.help.helper", "custom"))
	span.SetAttributes(attribute.String("task.help.request_id", gorequest.GetRequestIDContext(ctx)))

	span.SetAttributes(attribute.String("task.single.info", gojson.JsonEncodeNoError(task)))

	// 任务回调函数
	if executionCallback != nil {

		// 执行
		err = executionCallback(ctx, &task)
		if err != nil {
			span.RecordError(err, trace.WithStackTrace(true))
			span.SetStatus(codes.Error, err.Error())
		}

	}

	span.End() // 结束OpenTelemetry链路追踪
	return
}

// DbRunSingleTask 运行单个任务
// ctx 链路追踪的上下文
// message 任务json编码
// executionCallback 执行任务回调函数 返回 runCode=状态 runDesc=描述
// updateCallback 执行更新回调函数
func (c *PubSubClient) DbRunSingleTask(ctx context.Context, message string, executionCallback func(ctx context.Context, task *GormModelTask) (runCode int, runDesc string), updateCallback func(ctx context.Context, task *GormModelTask, result *TaskHelperRunSingleTaskResponse)) {

	// 解析任务
	var task GormModelTask
	err := gojson.Unmarshal([]byte(message), &task)
	if err != nil {
		return
	}

	// 启动OpenTelemetry链路追踪
	ctx, span := NewTraceStartSpan(ctx, "DbRunSingleTask "+task.Type+" "+task.CustomID)

	span.SetAttributes(attribute.String("task.help.helper", "db"))
	span.SetAttributes(attribute.String("task.help.request_id", gorequest.GetRequestIDContext(ctx)))

	span.SetAttributes(attribute.String("task.single.info", gojson.JsonEncodeNoError(task)))

	// 任务回调函数
	if executionCallback != nil {

		// 需要返回的结构
		result := TaskHelperRunSingleTaskResponse{
			TraceID:   span.SpanContext().TraceID().String(),
			SpanID:    span.SpanContext().SpanID().String(),
			RequestID: gorequest.GetRequestIDContext(ctx),
		}

		// 执行
		result.RunCode, result.RunDesc = executionCallback(ctx, &task)
		if result.RunCode == CodeAbnormal {
			span.SetStatus(codes.Error, result.RunDesc)
		}
		if result.RunCode == CodeSuccess {
			span.SetStatus(codes.Ok, result.RunDesc)
		}
		if result.RunCode == CodeError {
			span.RecordError(fmt.Errorf(result.RunDesc), trace.WithStackTrace(true))
			span.SetStatus(codes.Error, result.RunDesc)
		}

		// 运行编号
		result.RunID = result.TraceID
		if result.RunID == "" {
			result.RunID = result.RequestID
			if result.RunID == "" {
				span.RecordError(fmt.Errorf("上下文没有运行编号"), trace.WithStackTrace(true))
				span.SetStatus(codes.Error, "上下文没有运行编号")

				span.End() // 结束OpenTelemetry链路追踪
				return
			}
		}

		// OpenTelemetry链路追踪
		span.SetAttributes(attribute.String("task.run.id", result.RunID))
		span.SetAttributes(attribute.Int("task.run.code", result.RunCode))
		span.SetAttributes(attribute.String("task.run.desc", result.RunDesc))

		// 执行更新回调函数
		if updateCallback != nil {
			updateCallback(ctx, &task, &result)
		}

	}

	span.End() // 结束OpenTelemetry链路追踪
	return
}
