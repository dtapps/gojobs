package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gorequest"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"strings"
)

type TaskHelper struct {
	traceIsFilter bool            // [过滤]链路追踪是否过滤
	newCtx        context.Context // [启动]上下文
	newSpan       trace.Span      // [启动]链路追踪
	listCtx       context.Context // [列表]上下文
	listSpan      trace.Span      // [列表]链路追踪
	//listTaskType  string          // [列表]任务类型
	//listTaskList  []GormModelTask // [列表]任务列表
	filterCtx  context.Context // [过滤]上下文
	filterSpan trace.Span      // [过滤]链路追踪
	runCtx     context.Context // [运行]上下文
	runSpan    trace.Span      // [运行]链路追踪
}

func NewTaskHelper(ctx context.Context, traceIsFilter bool) *TaskHelper {
	th := &TaskHelper{
		traceIsFilter: traceIsFilter,
	}

	// 启动OpenTelemetry链路追踪
	th.newCtx, th.newSpan = TraceStartSpan(ctx, "NewTaskHelper")

	return th
}

// QueryTaskList 通过回调函数获取任务列表
// callback 回调函数 返回 任务列表
// newTaskLists 新的任务列表
// isContinue 是否继续
func (th *TaskHelper) QueryTaskList(taskType string, callback func(ctx context.Context) (taskLists []GormModelTask)) (newTaskLists []GormModelTask, isContinue bool) {

	// 启动OpenTelemetry链路追踪
	th.listCtx, th.listSpan = TraceStartSpan(th.newCtx, "QueryTaskList "+taskType)

	// 任务类型不能为空
	if taskType != "" {

		if callback != nil {

			// 执行任务列表回调函数
			taskLists := callback(th.listCtx)
			if taskLists == nil {
				if th.traceIsFilter {
					TraceSetAttributes(th.newCtx, attribute.String("is_filter", "true"))
					TraceSetAttributes(th.listCtx, attribute.String("is_filter", "true"))
				}
				return
			}

			// 判断任务类型是否一致
			for _, vTask := range taskLists {
				if vTask.Type == taskType {
					newTaskLists = append(newTaskLists, vTask)
				}
			}

		}

	}

	// 没有任务需要执行
	if len(newTaskLists) <= 0 {
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
	TraceSetAttributes(th.listCtx, attribute.String("task.list.type", taskType))
	TraceSetAttributes(th.listCtx, attribute.Int("task.list.count", len(newTaskLists)))

	return newTaskLists, true
}

// FilterTaskList 过滤任务列表
// isMandatoryIp 强制当前ip
// specifyIp 指定Ip
// taskType 任务类型
// taskLists 任务类型
// taskLists 过滤前的数据
// newTaskLists 过滤后的数据
// isContinue 是否继续
func (th *TaskHelper) FilterTaskList(isMandatoryIp bool, specifyIp string, taskType string, taskLists []GormModelTask) (newTaskLists []GormModelTask, isContinue bool) {

	// 启动OpenTelemetry链路追踪
	th.filterCtx, th.filterSpan = TraceStartSpan(th.listCtx, "FilterTaskList "+taskType)

	if specifyIp != "" {

		// 解析指定IP
		specifyIp = gorequest.IpIs(specifyIp)

		for _, vTask := range taskLists {

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

	}

	// 没有任务需要执行
	if len(newTaskLists) <= 0 {
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
	TraceSetAttributes(th.filterCtx, attribute.String("task.filter.type", taskType))
	TraceSetAttributes(th.filterCtx, attribute.Int("task.filter.count", len(newTaskLists)))

	return newTaskLists, true
}

// RunTask 执行任务
// task 任务
// executionCallback 执行任务回调函数 返回 状态和描述
// updateCallback 更新回调函数
func (th *TaskHelper) RunTask(task GormModelTask, executionCallback func(ctx context.Context, task GormModelTask) (runCode int, runDesc string), updateCallback func(ctx context.Context, task GormModelTask, runCode int, runDesc string)) {

	// 启动OpenTelemetry链路追踪
	th.runCtx, th.runSpan = TraceStartSpan(th.filterCtx, "RunTask "+task.Type+" "+task.CustomID)

	if executionCallback != nil {

		// 执行任务回调函数
		runCode, runDesc := executionCallback(th.runCtx, task)
		if runCode == CodeAbnormal {
			TraceSetStatus(th.runCtx, codes.Error, runDesc)
		}
		if runCode == CodeSuccess {
			TraceSetStatus(th.runCtx, codes.Ok, runDesc)
		}
		if runCode == CodeError {
			TraceRecordError(th.runCtx, fmt.Errorf(runDesc))
			TraceSetStatus(th.runCtx, codes.Error, runDesc)
		}

		// 运行编号
		runID := TraceGetTraceID(th.runCtx)
		if runID == "" {
			runID = gorequest.GetRequestIDContext(th.runCtx)
			if runID == "" {
				TraceRecordError(th.runCtx, fmt.Errorf("上下文没有运行编号"))
				TraceSetStatus(th.runCtx, codes.Error, "上下文没有运行编号")

				// 停止OpenTelemetry链路追踪
				TraceEndSpan(th.runSpan)
				TraceEndSpan(th.filterSpan)
				TraceEndSpan(th.listSpan)
				TraceEndSpan(th.newSpan)
				return
			}
		}

		// OpenTelemetry链路追踪
		TraceSetAttributes(th.runCtx, attribute.Int64("task.info.id", int64(task.ID)))
		TraceSetAttributes(th.runCtx, attribute.String("task.info.status", task.Status))
		TraceSetAttributes(th.runCtx, attribute.String("task.info.params", task.Params))
		TraceSetAttributes(th.runCtx, attribute.Int64("task.info.number", task.Number))
		TraceSetAttributes(th.runCtx, attribute.Int64("task.info.max_number", task.MaxNumber))
		TraceSetAttributes(th.runCtx, attribute.String("task.info.custom_id", task.CustomID))
		TraceSetAttributes(th.runCtx, attribute.Int64("task.info.custom_sequence", task.CustomSequence))
		TraceSetAttributes(th.runCtx, attribute.String("task.info.type", task.Type))
		TraceSetAttributes(th.runCtx, attribute.String("task.info.type_name", task.TypeName))
		TraceSetAttributes(th.runCtx, attribute.String("task.run.id", runID))
		TraceSetAttributes(th.runCtx, attribute.Int("task.run.code", runCode))
		TraceSetAttributes(th.runCtx, attribute.String("task.run.desc", runDesc))

		// 执行更新回调函数
		if updateCallback != nil {
			updateCallback(th.runCtx, task, runCode, runDesc)
		}

	}

	// 停止OpenTelemetry链路追踪
	TraceEndSpan(th.runSpan)
	TraceEndSpan(th.filterSpan)
	TraceEndSpan(th.listSpan)
	TraceEndSpan(th.newSpan)
	return
}
