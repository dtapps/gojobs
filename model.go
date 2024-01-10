package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
)

// 创建模型
func (c *Client) autoMigrateTask(ctx context.Context) {
	err := c.gormClient.AutoMigrate(
		&jobs_gorm_model.Task{},
		&jobs_gorm_model.TaskLog{},
	)
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("创建模型：%s", err))
		}
	}
}

// GormTaskLogDelete 删除
func (c *Client) GormTaskLogDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogInDelete 删除任务运行
func (c *Client) GormTaskLogInDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.Where("task_result_status = ?", TASK_IN).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogSuccessDelete 删除任务完成
func (c *Client) GormTaskLogSuccessDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.Where("task_result_status = ?", TASK_SUCCESS).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogErrorDelete 删除任务异常
func (c *Client) GormTaskLogErrorDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.Where("task_result_status = ?", TASK_ERROR).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogTimeoutDelete 删除任务超时
func (c *Client) GormTaskLogTimeoutDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.Where("task_result_status = ?", TASK_TIMEOUT).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogWaitDelete 删除任务等待
func (c *Client) GormTaskLogWaitDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.Where("task_result_status = ?", TASK_WAIT).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// TaskLogRecord 记录
func (c *Client) TaskLogRecord(ctx context.Context, task jobs_gorm_model.Task, runId string, taskResultCode int, taskResultDesc string) {
	c.GormTaskLogRecord(ctx, task, runId, taskResultCode, taskResultDesc)
}

// GormTaskLogRecord 记录
func (c *Client) GormTaskLogRecord(ctx context.Context, task jobs_gorm_model.Task, runId string, taskResultCode int, taskResultDesc string) {

	taskLog := jobs_gorm_model.TaskLog{
		TaskID:         task.ID,
		TaskRunID:      runId,
		TaskResultCode: taskResultCode,
		TaskResultDesc: taskResultDesc,

		SystemHostName:  c.config.systemHostname,
		SystemInsideIP:  c.config.systemInsideIP,
		SystemOutsideIP: c.config.systemOutsideIP,
		SystemOs:        c.config.systemOs,
		SystemArch:      c.config.systemKernel,
		SystemUpTime:    c.config.systemUpTime,
		SystemBootTime:  c.config.systemBootTime,
		GoVersion:       c.config.goVersion,
		SdkVersion:      c.config.sdkVersion,
		SystemVersion:   c.config.sdkVersion,
		CpuCores:        c.config.cpuCores,
		CpuModelName:    c.config.cpuModelName,
		CpuMhz:          c.config.cpuMhz,
	}
	err := c.gormClient.Create(&taskLog).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录失败：%s", err))
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("记录数据：%+v", taskLog))
		}
	}

}
