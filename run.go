package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
)

// Run 运行
func (c *Client) Run(ctx context.Context, info jobs_gorm_model.Task, status int, result string) {

	runId := gotrace_id.GetTraceIdContext(ctx)
	if runId == "" {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run]：%s", "上下文没有跟踪编号")
		return
	}

	switch status {
	case 0:
		err := c.EditTask(c.db.gormClient.Db, info.Id).
			Select("run_id", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				RunId:       runId,
				Result:      result,
				NextRunTime: gotime.Current().AfterSeconds(info.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.0]：%s", err.Error())
		}
		return
	case CodeSuccess:
		// 执行成功
		err := c.EditTask(c.db.gormClient.Db, info.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行成功",
				Number:      info.Number + 1,
				RunId:       runId,
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      result,
				NextRunTime: gotime.Current().AfterSeconds(info.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.CodeSuccess]：%s", err.Error())
		}
	case CodeEnd:
		// 执行成功、提前结束
		err := c.EditTask(c.db.gormClient.Db, info.Id).
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
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.CodeEnd]：%s", err.Error())
		}
	case CodeError:
		// 执行失败
		err := c.EditTask(c.db.gormClient.Db, info.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行失败",
				Number:      info.Number + 1,
				RunId:       runId,
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      result,
				NextRunTime: gotime.Current().AfterSeconds(info.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.CodeError]：%s", err.Error())
		}
	}

	if info.MaxNumber != 0 {
		if info.Number+1 >= info.MaxNumber {
			// 关闭执行
			err := c.EditTask(c.db.gormClient.Db, info.Id).
				Select("status").
				Updates(jobs_gorm_model.Task{
					Status: TASK_TIMEOUT,
				}).Error
			if err != nil {
				c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.TASK_TIMEOUT]：%s", err.Error())
			}
		}
	}
	return
}
