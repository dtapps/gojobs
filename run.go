package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gostring"
	"go.dtapp.net/gotime"
)

// Run 运行
func (c *Client) Run(ctx context.Context, info jobs_gorm_model.Task, status int, result string) {
	// 请求函数记录
	err := c.db.gormClient.Db.Create(&jobs_gorm_model.TaskLog{
		TaskId:     info.Id,
		StatusCode: status,
		Desc:       result,
		Version:    c.config.sdkVersion,
	}).Error
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.Create]：%s", err.Error())
	}
	if status == 0 {
		err = c.EditTask(c.db.gormClient.Db, info.Id).
			Select("run_id", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				RunId:       gostring.GetUuId(),
				Result:      result,
				NextRunTime: gotime.Current().AfterSeconds(info.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.0.EditTask]：%s", err.Error())
		}
		return
	}
	// 任务
	if status == CodeSuccess {
		// 执行成功
		err = c.EditTask(c.db.gormClient.Db, info.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行成功",
				Number:      info.Number + 1,
				RunId:       gostring.GetUuId(),
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      result,
				NextRunTime: gotime.Current().AfterSeconds(info.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.CodeSuccess.EditTask]：%s", err.Error())
		}
	}
	if status == CodeEnd {
		// 执行成功、提前结束
		err = c.EditTask(c.db.gormClient.Db, info.Id).
			Select("status", "status_desc", "number", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				Status:      TASK_SUCCESS,
				StatusDesc:  "结束执行",
				Number:      info.Number + 1,
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      result,
				NextRunTime: gotime.Current().Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.CodeEnd.EditTask]：%s", err.Error())
		}
	}
	if status == CodeError {
		// 执行失败
		err = c.EditTask(c.db.gormClient.Db, info.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行失败",
				Number:      info.Number + 1,
				RunId:       gostring.GetUuId(),
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      result,
				NextRunTime: gotime.Current().AfterSeconds(info.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.CodeError.EditTask]：%s", err.Error())
		}
	}
	if info.MaxNumber != 0 {
		if info.Number+1 >= info.MaxNumber {
			// 关闭执行
			err = c.EditTask(c.db.gormClient.Db, info.Id).
				Select("status").
				Updates(jobs_gorm_model.Task{
					Status: TASK_TIMEOUT,
				}).Error
			if err != nil {
				c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.TASK_TIMEOUT.EditTask]：%s", err.Error())
			}
		}
	}
}

// RunAddLog 任务执行日志
func (c *Client) RunAddLog(ctx context.Context, id uint, runId string) error {
	return c.db.gormClient.Db.Create(&jobs_gorm_model.TaskLogRun{
		TaskId:     id,
		RunId:      runId,
		InsideIp:   c.config.systemInsideIp,
		OutsideIp:  c.config.systemOutsideIp,
		Os:         c.config.systemOs,
		Arch:       c.config.systemArch,
		Gomaxprocs: c.config.systemCpuQuantity,
		GoVersion:  c.config.goVersion,
		SdkVersion: c.config.sdkVersion,
		MacAddrs:   c.config.systemMacAddrS,
	}).Error
}
