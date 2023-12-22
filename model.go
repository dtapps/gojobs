package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gojson"
)

// 创建模型
func (c *Client) autoMigrateTask(ctx context.Context) {
	err := c.gormClient.GetDb().AutoMigrate(&jobs_gorm_model.Task{})
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("创建模型：%s", err))
		}
	}
}

// TaskLog 任务日志模型
type TaskLog struct {
	TaskId           uint   `json:"task_id"`            // 任务编号
	TaskRunId        string `json:"task_run_id"`        // 执行编号
	TaskResultStatus string `json:"task_result_status"` // 执行状态
	TaskResultCode   int    `json:"task_result_code"`   // 执行状态码
	TaskResultDesc   string `json:"task_result_desc"`   // 执行结果
	SystemInsideIp   string `json:"system_inside_ip"`   // 内网ip
	SystemOutsideIp  string `json:"system_outside_ip"`  // 外网ip
}

// TaskLogRecord 记录
func (c *Client) TaskLogRecord(ctx context.Context, task jobs_gorm_model.Task, runId string, taskResultCode int, taskResultDesc string) {
	if c.runSlog.status {
		taskLog := TaskLog{
			TaskId:          task.Id,
			TaskRunId:       runId,
			TaskResultCode:  taskResultCode,
			TaskResultDesc:  taskResultDesc,
			SystemInsideIp:  c.config.systemInsideIp,
			SystemOutsideIp: c.config.systemOutsideIp,
		}
		c.runSlog.client.WithTraceId(ctx).Info(gojson.JsonEncodeNoError(taskLog))
	}
}
