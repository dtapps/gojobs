package gojobs

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"go.dtapp.net/gojson"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
	"strings"
	"time"
)

type TaskHelper struct {
	cfg *taskHelperConfig // 配置

	taskType string          // [任务]类型
	taskList []GormModelTask // [任务]列表

	newCtx            context.Context // [启动]上下文
	newSpan           trace.Span      // [启动]链路追踪
	listCtx           context.Context // [列表]上下文
	listSpan          trace.Span      // [列表]链路追踪
	filterCtx         context.Context // [过滤]上下文
	filterSpan        trace.Span      // [过滤]链路追踪
	runMultipleStatus bool            // [运行多个]状态
	runMultipleCtx    context.Context // [运行多个]上下文
	runMultipleSpan   trace.Span      // [运行多个]链路追踪
	runSingleCtx      context.Context // [运行单个]上下文
	runSingleSpan     trace.Span      // [运行单个]链路追踪
}

// NewTaskHelper 任务帮助
// ctx 链路追踪的上下文
// taskType 任务类型
// logIsDebug 日志是否启动
// traceIsFilter 链路追踪是否过滤
func NewTaskHelper(ctx context.Context, taskType string, opts ...TaskHelperOption) (*TaskHelper, error) {
	th := &TaskHelper{
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

	th.newSpan.SetAttributes(attribute.Bool("task.cfg.is_debug", th.cfg.logIsDebug))
	th.newSpan.SetAttributes(attribute.Bool("task.cfg.is_filter", th.cfg.traceIsFilter))
	th.newSpan.SetAttributes(attribute.String("task.cfg.filter_name", th.cfg.traceIsFilterKeyName))
	th.newSpan.SetAttributes(attribute.String("task.cfg.filter_value", th.cfg.traceIsFilterKeyValue))

	return th, nil
}

// QueryTaskList 通过回调函数获取任务列表
// isRunCallback 任务列表回调函数 返回 是否使用 任务列表
// listCallback 任务回调函数 返回 任务列表
// newTaskLists 新的任务列表
// isContinue 是否继续
func (th *TaskHelper) QueryTaskList(isRunCallback func(ctx context.Context, keyName string) (isUse bool, result *redis.StringCmd), listCallback func(ctx context.Context, taskType string) []GormModelTask) (isContinue bool) {

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
					//th.listSpan.RecordError( err, trace.WithStackTrace(true))
					th.listSpan.SetStatus(codes.Error, err.Error())

					if th.cfg.logIsDebug {
						slog.DebugContext(th.listCtx, "查询redis的key不存在，根据设置，无法继续运行", slog.String("key", GetRedisKeyName(th.taskType)), slog.String("err", isRunResult.Err().Error()))
					}

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
				//th.listSpan.RecordError( err, trace.WithStackTrace(true))
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
				//th.listSpan.RecordError(err, trace.WithStackTrace(true))
				th.listSpan.SetStatus(codes.Error, err.Error())

				if th.cfg.logIsDebug {
					slog.DebugContext(th.listCtx, "QueryTaskList 查询redis的key内容为空，根据配置，无法继续运行", slog.String("val", isRunResult.Val()))
				}

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
		taskLists := listCallback(th.listCtx, th.taskType)

		// 判断任务类型是否一致
		for _, vTask := range taskLists {
			if vTask.Type == th.taskType {
				th.taskList = append(th.taskList, vTask)
			}
		}

	}

	// 没有任务需要执行
	if len(th.taskList) <= 0 {
		if th.cfg.logIsDebug {
			slog.InfoContext(th.listCtx, "QueryTaskList 没有任务需要执行")
		}

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

// FilterTaskList 过滤任务列表
// isMandatoryIp 强制当前ip
// specifyIp 指定Ip
// isContinue 是否继续
func (th *TaskHelper) FilterTaskList(isMandatoryIp bool, specifyIp string) (isContinue bool) {

	// 启动OpenTelemetry链路追踪
	th.filterCtx, th.filterSpan = NewTraceStartSpan(th.listCtx, "FilterTaskList")

	if th.cfg.logIsDebug {
		slog.DebugContext(th.filterCtx, "FilterTaskList 过滤任务列表", slog.Bool("isMandatoryIp", isMandatoryIp), slog.String("specifyIp", specifyIp))
	}

	if specifyIp != "" {

		// 新的任务列表
		var newTaskLists []GormModelTask

		// 解析指定IP
		specifyIp = gorequest.IpIs(specifyIp)

		// 循环判断 过滤指定IP
		for _, vTask := range th.taskList {

			vTask.SpecifyIP = gorequest.IpIs(vTask.SpecifyIP)

			// 强制只能是当前的IP
			if isMandatoryIp {
				// 进入强制性IP
				if vTask.SpecifyIP == specifyIp {
					// 进入强制性IP，可添加任务
					newTaskLists = append(newTaskLists, vTask)
					continue
				}
			}

			if vTask.SpecifyIP == "" {
				// 任务指定IP为空，可添加任务
				newTaskLists = append(newTaskLists, vTask)
				continue
			} else if vTask.SpecifyIP == SpecifyIpNull {
				// 任务指定Ip无限制，可添加任务
				newTaskLists = append(newTaskLists, vTask)
				continue
			} else {
				// 判断是否包含该IP
				specifyIpFind := strings.Contains(vTask.SpecifyIP, ",")
				if specifyIpFind {
					// 进入强制性多IP
					// 分割字符串
					parts := strings.Split(vTask.SpecifyIP, ",")
					for _, vv := range parts {
						if vv == specifyIp {
							// 进入强制性多IP，可添加任务
							newTaskLists = append(newTaskLists, vTask)
							continue
						}
					}
				} else {
					// 进入强制性单IP
					if vTask.SpecifyIP == specifyIp {
						// 进入强制性单IP，可添加任务
						newTaskLists = append(newTaskLists, vTask)
						continue
					}
				}
			}
		}

		// 设置任务列表
		th.taskList = newTaskLists
	}

	// 没有任务需要执行
	if len(th.taskList) <= 0 {
		if th.cfg.logIsDebug {
			slog.InfoContext(th.filterCtx, "FilterTaskList 没有任务需要执行")
		}

		if th.cfg.traceIsFilter {
			th.filterSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
			th.listSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
			th.newSpan.SetAttributes(attribute.String(th.cfg.traceIsFilterKeyName, th.cfg.traceIsFilterKeyValue))
		}

		// 停止OpenTelemetry链路追踪
		th.filterSpan.End()
		th.listSpan.End()
		th.newSpan.End()
		return
	}

	// OpenTelemetry链路追踪
	th.filterSpan.SetAttributes(attribute.Int("task.filter.count", len(th.taskList)))

	return true
}

// GetTaskList 获取任务列表
func (th *TaskHelper) GetTaskList() []GormModelTask {
	return th.taskList
}

// RunMultipleTask 运行多个任务
// executionCallback 执行任务回调函数 返回 runCode=状态 runDesc=描述
// updateCallback 执行更新回调函数
func (th *TaskHelper) RunMultipleTask(wait int64, executionCallback func(ctx context.Context, task GormModelTask) (runCode int, runDesc string), updateCallback func(ctx context.Context, task GormModelTask, result RunSingleTaskResponse)) {

	// 启动OpenTelemetry链路追踪
	th.runMultipleStatus = true
	th.runMultipleCtx, th.runMultipleSpan = NewTraceStartSpan(th.filterCtx, "RunMultipleTask")

	if th.cfg.logIsDebug {
		slog.DebugContext(th.runMultipleCtx, "RunMultipleTask 运行多个任务", slog.Int64("wait", wait))
	}

	for _, vTask := range th.taskList {

		// 运行单个任务
		th.RunSingleTask(vTask, executionCallback, updateCallback)

		// 等待 wait 秒
		if wait > 0 {
			time.Sleep(time.Duration(wait) * time.Second)
		}

	}

	// 停止OpenTelemetry链路追踪
	th.EndRunTaskList()
	return
}

type RunSingleTaskResponse struct {
	RunID   string // 运行编号
	RunCode int    // 运行状态
	RunDesc string // 运行描述

	TraceID   string // 追踪编号
	SpanID    string // 跨度编号
	RequestID string // 请求编号
}

// RunSingleTask 运行单个任务
// task 任务
// executionCallback 执行任务回调函数 返回 runCode=状态 runDesc=描述
// updateCallback 执行更新回调函数
func (th *TaskHelper) RunSingleTask(task GormModelTask, executionCallback func(ctx context.Context, task GormModelTask) (runCode int, runDesc string), updateCallback func(ctx context.Context, task GormModelTask, result RunSingleTaskResponse)) {

	// 启动OpenTelemetry链路追踪
	if th.runMultipleStatus {
		th.runSingleCtx, th.runSingleSpan = NewTraceStartSpan(th.runMultipleCtx, "RunSingleTask "+task.CustomID)
	} else {
		th.runSingleCtx, th.runSingleSpan = NewTraceStartSpan(th.filterCtx, "RunSingleTask "+th.taskType+" "+task.CustomID)
	}

	if th.cfg.logIsDebug {
		slog.DebugContext(th.runSingleCtx, "RunSingleTask 运行单个任务", slog.String("task", gojson.JsonEncodeNoError(task)))
	}

	// 任务回调函数
	if executionCallback != nil {

		// 需要返回的结构
		result := RunSingleTaskResponse{
			TraceID:   th.runSingleSpan.SpanContext().TraceID().String(),
			SpanID:    th.runSingleSpan.SpanContext().SpanID().String(),
			RequestID: gorequest.GetRequestIDContext(th.runSingleCtx),
		}

		// 执行
		result.RunCode, result.RunDesc = executionCallback(th.runSingleCtx, task)
		if result.RunCode == CodeAbnormal {
			th.runSingleSpan.SetStatus(codes.Error, result.RunDesc)
		}
		if result.RunCode == CodeSuccess {
			th.runSingleSpan.SetStatus(codes.Ok, result.RunDesc)
		}
		if result.RunCode == CodeError {
			th.runSingleSpan.RecordError(fmt.Errorf(result.RunDesc), trace.WithStackTrace(true))
			th.runSingleSpan.SetStatus(codes.Error, result.RunDesc)
		}

		// 运行编号
		result.RunID = result.TraceID
		if result.RunID == "" {
			result.RunID = result.RequestID
			if result.RunID == "" {
				th.runSingleSpan.RecordError(fmt.Errorf("上下文没有运行编号"), trace.WithStackTrace(true))
				th.runSingleSpan.SetStatus(codes.Error, "上下文没有运行编号")

				if th.cfg.logIsDebug {
					slog.ErrorContext(th.listCtx, "RunSingleTask 上下文没有运行编号")
				}

				// 停止OpenTelemetry链路追踪
				th.runSingleSpan.End()
				return
			}
		}

		// OpenTelemetry链路追踪
		th.runSingleSpan.SetAttributes(attribute.Int64("task.info.id", int64(task.ID)))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.status", task.Status))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.params", task.Params))
		th.runSingleSpan.SetAttributes(attribute.Int64("task.info.number", task.Number))
		th.runSingleSpan.SetAttributes(attribute.Int64("task.info.max_number", task.MaxNumber))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.custom_id", task.CustomID))
		th.runSingleSpan.SetAttributes(attribute.Int64("task.info.custom_sequence", task.CustomSequence))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.type", task.Type))
		th.runSingleSpan.SetAttributes(attribute.String("task.info.type_name", task.TypeName))
		th.runSingleSpan.SetAttributes(attribute.String("task.run.id", result.RunID))
		th.runSingleSpan.SetAttributes(attribute.Int("task.run.code", result.RunCode))
		th.runSingleSpan.SetAttributes(attribute.String("task.run.desc", result.RunDesc))

		// 执行更新回调函数
		if updateCallback != nil {
			updateCallback(th.runSingleCtx, task, result)
		}

	}

	// 停止OpenTelemetry链路追踪
	th.runSingleSpan.End()
	return
}

// EndRunTaskList 结束运行任务列表并停止OpenTelemetry链路追踪
func (th *TaskHelper) EndRunTaskList() {
	if th.runMultipleStatus {
		th.runMultipleSpan.End()
	}
	th.filterSpan.End()
	th.listSpan.End()
	th.newSpan.End()
	return
}
