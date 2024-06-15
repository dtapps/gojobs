package gojobs

import (
	"context"
	"fmt"
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
	logIsDebug bool // [日志]日志是否启动

	traceIsFilter bool            // [过滤]链路追踪是否过滤
	taskType      string          // [任务]类型
	taskList      []GormModelTask // [任务]列表

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
func NewTaskHelper(ctx context.Context, taskType string, logIsDebug bool, traceIsFilter bool) (*TaskHelper, error) {
	th := &TaskHelper{
		logIsDebug:    logIsDebug,
		traceIsFilter: traceIsFilter,
		taskType:      taskType,
	}

	if th.taskType == "" {
		return th, fmt.Errorf("请检查任务类型参数")
	}

	// 启动OpenTelemetry链路追踪
	th.newCtx, th.newSpan = NewTraceStartSpan(ctx, "NewTaskHelper")

	TraceSetAttributes(th.newCtx, attribute.String("task.new.type", th.taskType))
	TraceSetAttributes(th.newCtx, attribute.Bool("task.new.is_debug", th.logIsDebug))
	TraceSetAttributes(th.newCtx, attribute.Bool("task.new.is_filter", th.traceIsFilter))

	return th, nil
}

// QueryTaskList 通过回调函数获取任务列表
// callback 回调函数 返回 任务列表
// newTaskLists 新的任务列表
// isContinue 是否继续
func (th *TaskHelper) QueryTaskList(callback func(ctx context.Context, taskType string) (taskLists []GormModelTask)) (isContinue bool) {

	// 启动OpenTelemetry链路追踪
	th.listCtx, th.listSpan = NewTraceStartSpan(th.newCtx, "QueryTaskList "+th.taskType)

	if callback != nil {

		// 执行任务列表回调函数
		taskLists := callback(th.listCtx, th.taskType)
		if taskLists == nil {
			if th.traceIsFilter {
				TraceSetAttributes(th.newCtx, attribute.String("is_filter", "true"))
				TraceSetAttributes(th.listCtx, attribute.String("is_filter", "true"))
			}
			return
		}

		// 判断任务类型是否一致
		for _, vTask := range taskLists {
			if vTask.Type == th.taskType {
				th.taskList = append(th.taskList, vTask)
			}
		}

	}

	// 没有任务需要执行
	if len(th.taskList) <= 0 {
		if th.logIsDebug {
			slog.ErrorContext(th.listCtx, "QueryTaskList 没有任务需要执行", slog.Int("taskLists", len(th.taskList)))
		}
		if th.traceIsFilter {
			TraceSetAttributes(th.newCtx, attribute.String("is_filter", "true"))
			TraceSetAttributes(th.listCtx, attribute.String("is_filter", "true"))
		}
		// 停止OpenTelemetry链路追踪
		TraceEndSpan(th.listSpan)
		TraceEndSpan(th.newSpan)
		return
	}

	// OpenTelemetry链路追踪
	TraceSetAttributes(th.listCtx, attribute.Int("task.list.count", len(th.taskList)))

	return true
}

// FilterTaskList 过滤任务列表
// isMandatoryIp 强制当前ip
// specifyIp 指定Ip
// isContinue 是否继续
func (th *TaskHelper) FilterTaskList(isMandatoryIp bool, specifyIp string) (isContinue bool) {

	// 启动OpenTelemetry链路追踪
	th.filterCtx, th.filterSpan = NewTraceStartSpan(th.listCtx, "FilterTaskList "+th.taskType)

	if th.logIsDebug {
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
		if th.logIsDebug {
			slog.ErrorContext(th.filterCtx, "FilterTaskList 没有任务需要执行", slog.Bool("isMandatoryIp", isMandatoryIp), slog.String("specifyIp", specifyIp))
		}
		if th.traceIsFilter {
			TraceSetAttributes(th.newCtx, attribute.String("is_filter", "true"))
			TraceSetAttributes(th.listCtx, attribute.String("is_filter", "true"))
			TraceSetAttributes(th.filterCtx, attribute.String("is_filter", "true"))
		}
		// 停止OpenTelemetry链路追踪
		TraceEndSpan(th.filterSpan)
		TraceEndSpan(th.listSpan)
		TraceEndSpan(th.newSpan)
		return
	}

	// OpenTelemetry链路追踪
	TraceSetAttributes(th.filterCtx, attribute.Int("task.filter.count", len(th.taskList)))

	return true
}

// GetTaskList 获取任务列表
func (th *TaskHelper) GetTaskList() []GormModelTask {
	return th.taskList
}

// RunMultipleTask 运行多个任务
// executionCallback 执行任务回调函数 返回 状态和描述
// updateCallback 更新回调函数
func (th *TaskHelper) RunMultipleTask(wait int64, executionCallback func(ctx context.Context, task GormModelTask) (runCode int, runDesc string), updateCallback func(ctx context.Context, task GormModelTask, runCode int, runDesc string)) {

	// 启动OpenTelemetry链路追踪
	th.runMultipleStatus = true
	th.runMultipleCtx, th.runMultipleSpan = NewTraceStartSpan(th.filterCtx, "RunMultipleTask "+th.taskType)

	if th.logIsDebug {
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

// RunSingleTask 运行单个任务
// task 任务
// executionCallback 执行任务回调函数 返回 状态和描述
// updateCallback 更新回调函数
func (th *TaskHelper) RunSingleTask(task GormModelTask, executionCallback func(ctx context.Context, task GormModelTask) (runCode int, runDesc string), updateCallback func(ctx context.Context, task GormModelTask, runCode int, runDesc string)) {

	// 启动OpenTelemetry链路追踪
	if th.runMultipleStatus {
		th.runSingleCtx, th.runSingleSpan = NewTraceStartSpan(th.runMultipleCtx, "RunSingleTask "+th.taskType+" "+task.CustomID)
	} else {
		th.runSingleCtx, th.runSingleSpan = NewTraceStartSpan(th.filterCtx, "RunSingleTask "+th.taskType+" "+task.CustomID)
	}

	if th.logIsDebug {
		slog.DebugContext(th.runSingleCtx, "RunSingleTask 运行单个任务", slog.String("task", gojson.JsonEncodeNoError(task)))
	}

	if executionCallback != nil {

		// 执行任务回调函数
		runCode, runDesc := executionCallback(th.runSingleCtx, task)
		if runCode == CodeAbnormal {
			TraceSetStatus(th.runSingleCtx, codes.Error, runDesc)
		}
		if runCode == CodeSuccess {
			TraceSetStatus(th.runSingleCtx, codes.Ok, runDesc)
		}
		if runCode == CodeError {
			TraceRecordError(th.runSingleCtx, fmt.Errorf(runDesc))
			TraceSetStatus(th.runSingleCtx, codes.Error, runDesc)
		}

		// 运行编号
		runID := TraceGetTraceID(th.runSingleCtx)
		if runID == "" {
			runID = gorequest.GetRequestIDContext(th.runSingleCtx)
			if runID == "" {
				TraceRecordError(th.runSingleCtx, fmt.Errorf("上下文没有运行编号"))
				TraceSetStatus(th.runSingleCtx, codes.Error, "上下文没有运行编号")

				// 停止OpenTelemetry链路追踪
				TraceEndSpan(th.runSingleSpan)
				return
			}
		}

		// OpenTelemetry链路追踪
		TraceSetAttributes(th.runSingleCtx, attribute.Int64("task.info.id", int64(task.ID)))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.info.status", task.Status))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.info.params", task.Params))
		TraceSetAttributes(th.runSingleCtx, attribute.Int64("task.info.number", task.Number))
		TraceSetAttributes(th.runSingleCtx, attribute.Int64("task.info.max_number", task.MaxNumber))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.info.custom_id", task.CustomID))
		TraceSetAttributes(th.runSingleCtx, attribute.Int64("task.info.custom_sequence", task.CustomSequence))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.info.type", task.Type))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.info.type_name", task.TypeName))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.run.id", runID))
		TraceSetAttributes(th.runSingleCtx, attribute.Int("task.run.code", runCode))
		TraceSetAttributes(th.runSingleCtx, attribute.String("task.run.desc", runDesc))

		// 执行更新回调函数
		if updateCallback != nil {
			updateCallback(th.runSingleCtx, task, runCode, runDesc)
		}

	}

	// 停止OpenTelemetry链路追踪
	TraceEndSpan(th.runSingleSpan)
	return
}

// EndRunTaskList 结束运行任务列表并停止OpenTelemetry链路追踪
func (th *TaskHelper) EndRunTaskList() {
	if th.runMultipleStatus {
		TraceEndSpan(th.runMultipleSpan)
	}
	TraceEndSpan(th.filterSpan)
	TraceEndSpan(th.listSpan)
	TraceEndSpan(th.newSpan)
	return
}
