package gojobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.dtapp.net/gojson"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"time"
)

type TaskCustomHelper struct {
	cfg *taskHelperConfig // 配置

	taskType string                     // [任务]类型
	taskList []TaskCustomHelperTaskList // [任务]列表

	newCtx            context.Context // [启动]上下文
	newSpan           trace.Span      // [启动]链路追踪
	listCtx           context.Context // [列表]上下文
	listSpan          trace.Span      // [列表]链路追踪
	runMultipleStatus bool            // [运行多个]状态
	runMultipleCtx    context.Context // [运行多个]上下文
	runMultipleSpan   trace.Span      // [运行多个]链路追踪
	runSingleCtx      context.Context // [运行单个]上下文
	runSingleSpan     trace.Span      // [运行单个]链路追踪
}

// NewTaskCustomHelper 任务帮助
// ctx 链路追踪的上下文
// taskType 任务类型
// logIsDebug 日志是否启动
// traceIsFilter 链路追踪是否过滤
func NewTaskCustomHelper(ctx context.Context, taskType string, opts ...TaskHelperOption) (*TaskCustomHelper, error) {
	th := &TaskCustomHelper{
		taskType: taskType,
	}

	if th.taskType == "" {
		return th, fmt.Errorf("请检查任务类型参数")
	}

	// 配置
	th.cfg = newTaskHelperConfig(opts)

	// 启动OpenTelemetry链路追踪
	th.newCtx, th.newSpan = NewTraceStartSpan(ctx, th.taskType)

	th.newSpan.SetAttributes(attribute.String("task.new.type", th.taskType))

	th.newSpan.SetAttributes(attribute.Bool("task.cfg.logIsDebug", th.cfg.logIsDebug))
	th.newSpan.SetAttributes(attribute.Bool("task.cfg.traceIsFilter", th.cfg.traceIsFilter))
	th.newSpan.SetAttributes(attribute.String("task.cfg.traceIsFilterKeyName", th.cfg.traceIsFilterKeyName))
	th.newSpan.SetAttributes(attribute.String("task.cfg.traceIsFilterKeyValue", th.cfg.traceIsFilterKeyValue))

	return th, nil
}

// QueryTaskList 通过回调函数获取任务列表
// isRunCallback 任务列表回调函数 返回 是否使用 任务列表
// listCallback 任务回调函数 返回 任务列表
// newTaskLists 新的任务列表
// isContinue 是否继续
func (th *TaskCustomHelper) QueryTaskList(isRunCallback func(ctx context.Context, keyName string) (isUse bool, result *redis.StringCmd), listCallback func(ctx context.Context, taskType string) []TaskCustomHelperTaskList) (isContinue bool) {

	// 启动OpenTelemetry链路追踪
	th.listCtx, th.listSpan = NewTraceStartSpan(th.newCtx, "QueryTaskList")

	// 任务列表回调函数
	if isRunCallback != nil {

		// 执行
		isRunUse, isRunResult := isRunCallback(th.listCtx, GetRedisKeyName(th.taskType))
		if isRunUse {
			if isRunResult.Err() != nil {
				if errors.Is(isRunResult.Err(), redis.Nil) {
					err := fmt.Errorf("查询redis的key不存在，根据设置，无法继续运行: %v", isRunResult.Err().Error())
					th.listSpan.SetStatus(codes.Error, err.Error())

					if th.cfg.logIsDebug {
						slog.DebugContext(th.listCtx, "查询redis的key不存在，根据设置，无法继续运行", slog.String("key", GetRedisKeyName(th.taskType)), slog.String("err", isRunResult.Err().Error()))
					}

					// 过滤
					if th.cfg.traceIsFilter {
						th.listSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
						th.newSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
					}

					// 停止OpenTelemetry链路追踪
					th.listSpan.End()
					th.newSpan.End()
					return
				}

				err := fmt.Errorf("查询redis的key异常，无法继续运行: %v", isRunResult.Err().Error())
				th.listSpan.SetStatus(codes.Error, err.Error())

				if th.cfg.logIsDebug {
					slog.DebugContext(th.listCtx, "QueryTaskList 查询redis的key异常，无法继续运行", slog.String("err", isRunResult.Err().Error()))
				}

				// 停止OpenTelemetry链路追踪
				th.listSpan.End()
				th.newSpan.End()
				return
			}
			if isRunResult.Val() == "" {
				err := fmt.Errorf("查询redis的key内容为空，根据配置，无法继续运行: %s", isRunResult.Val())
				th.listSpan.SetStatus(codes.Error, err.Error())

				if th.cfg.logIsDebug {
					slog.DebugContext(th.listCtx, "QueryTaskList 查询redis的key内容为空，根据配置，无法继续运行", slog.String("val", isRunResult.Val()))
				}

				// 过滤
				if th.cfg.traceIsFilter {
					th.listSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
					th.newSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
				}

				// 停止OpenTelemetry链路追踪
				th.listSpan.End()
				th.newSpan.End()
				return
			}
		}
	}

	// 任务列表回调函数
	if listCallback != nil {

		// 执行
		th.taskList = listCallback(th.listCtx, th.taskType)

	}

	// 没有任务需要执行
	if len(th.taskList) <= 0 {
		if th.cfg.logIsDebug {
			slog.InfoContext(th.listCtx, "QueryTaskList 没有任务需要执行")
		}

		// 过滤
		if th.cfg.traceIsFilter {
			th.listSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
			th.newSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
		}

		// 停止OpenTelemetry链路追踪
		th.listSpan.End()
		th.newSpan.End()
		return
	}

	// OpenTelemetry链路追踪
	th.listSpan.SetAttributes(attribute.Int("task.list.count", len(th.taskList)))

	return true
}

// GetTaskList 获取任务列表
func (th *TaskCustomHelper) GetTaskList() []TaskCustomHelperTaskList {
	return th.taskList
}

// RunMultipleTask 运行多个任务
// executionCallback 执行任务回调函数
func (th *TaskCustomHelper) RunMultipleTask(wait int64, executionCallback func(ctx context.Context, task TaskCustomHelperTaskList) (err error)) {

	// 启动OpenTelemetry链路追踪
	th.runMultipleStatus = true
	th.runMultipleCtx, th.runMultipleSpan = NewTraceStartSpan(th.listCtx, "RunMultipleTask")

	if th.cfg.logIsDebug {
		slog.DebugContext(th.runMultipleCtx, "RunMultipleTask 运行多个任务", slog.Int64("wait", wait))
	}

	for _, vTask := range th.taskList {

		// 运行单个任务
		th.RunSingleTask(vTask, executionCallback)

		// 等待 wait 秒
		if wait > 0 {
			time.Sleep(time.Duration(wait) * time.Second)
		}

	}

	// 停止OpenTelemetry链路追踪
	th.EndRunTaskList()
	return
}

// RunSingleTask 运行单个任务
// task 任务
// executionCallback 执行任务回调函数
func (th *TaskCustomHelper) RunSingleTask(task TaskCustomHelperTaskList, executionCallback func(ctx context.Context, task TaskCustomHelperTaskList) (err error)) {

	// 启动OpenTelemetry链路追踪
	if th.runMultipleStatus {
		th.runSingleCtx, th.runSingleSpan = NewTraceStartSpan(th.runMultipleCtx, "RunSingleTask "+task.CustomID)
	} else {
		th.runSingleCtx, th.runSingleSpan = NewTraceStartSpan(th.listCtx, "RunSingleTask "+task.CustomID)
	}

	if th.cfg.logIsDebug {
		slog.DebugContext(th.runSingleCtx, "RunSingleTask 运行单个任务", slog.String("task", gojson.JsonEncodeNoError(task)))
	}

	// 任务回调函数
	if executionCallback != nil {

		// 执行
		err := executionCallback(th.runSingleCtx, task)
		if err != nil {
			th.runSingleSpan.RecordError(err, trace.WithStackTrace(true))
			th.runSingleSpan.SetStatus(codes.Error, err.Error())
		}

		// OpenTelemetry链路追踪
		th.runSingleSpan.SetAttributes(attribute.String("task.info.id", task.TaskID))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.name", task.TaskName))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.params", task.TaskParams))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.custom_id", task.CustomID))

	}

	// 停止OpenTelemetry链路追踪
	th.runSingleSpan.End()
	return
}

// EndRunTaskList 结束运行任务列表并停止OpenTelemetry链路追踪
func (th *TaskCustomHelper) EndRunTaskList() {
	if th.runMultipleStatus {
		th.runMultipleSpan.End()
	}
	th.listSpan.End()
	th.newSpan.End()
	return
}
