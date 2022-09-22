package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Run 运行
func (c *Client) Run(ctx context.Context, task jobs_gorm_model.Task, taskResultCode int, taskResultDesc string) {

	runId := gotrace_id.GetTraceIdContext(ctx)
	if runId == "" {
		c.zapLog.WithTraceId(ctx).Sugar().Error("上下文没有跟踪编号")
		return
	}

	c.gormClient.GetDb().Create(&jobs_gorm_model.TaskLog{
		TaskId:          task.Id,
		TaskRunId:       runId,
		TaskResultCode:  taskResultCode,
		TaskResultDesc:  taskResultDesc,
		SystemHostName:  c.config.systemHostName,
		SystemInsideIp:  c.config.systemInsideIp,
		SystemOs:        c.config.systemOs,
		SystemArch:      c.config.systemArch,
		GoVersion:       c.config.goVersion,
		SdkVersion:      c.config.sdkVersion,
		SystemOutsideIp: c.config.systemOutsideIp,
	})
	if c.mongoConfig.stats {

		taskLog := TaskLog{
			LogId:   primitive.NewObjectID(),
			LogTime: primitive.NewDateTimeFromTime(gotime.Current().Time),
		}

		taskLog.Task.Id = task.Id
		taskLog.Task.RunId = runId
		taskLog.Task.ResultCode = taskResultCode
		taskLog.Task.ResultDesc = taskResultDesc
		taskLog.Task.ResultTime = dorm.NewBsonTimeCurrent()

		taskLog.System.HostName = c.config.systemHostName
		taskLog.System.InsideIp = c.config.systemInsideIp
		taskLog.System.OutsideIp = c.config.systemOutsideIp
		taskLog.System.Os = c.config.systemOs
		taskLog.System.Arch = c.config.systemArch

		taskLog.Version.Go = c.config.goVersion
		taskLog.Version.Sdk = c.config.sdkVersion

		c.mongoClient.Database(c.mongoConfig.databaseName).Collection(TaskLog{}.CollectionName()).InsertOne(ctx, taskLog)
	}

	switch taskResultCode {
	case 0:
		err := c.EditTask(c.gormClient.GetDb(), task.Id).
			Select("run_id", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				RunId:       runId,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().AfterSeconds(task.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("保存失败：%s", err.Error())
		}
		return
	case CodeSuccess:
		// 执行成功
		err := c.EditTask(c.gormClient.GetDb(), task.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行成功",
				Number:      task.Number + 1,
				RunId:       runId,
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().AfterSeconds(task.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("保存失败：%s", err.Error())
		}
	case CodeEnd:
		// 执行成功、提前结束
		err := c.EditTask(c.gormClient.GetDb(), task.Id).
			Select("status", "status_desc", "number", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				Status:      TASK_SUCCESS,
				StatusDesc:  "结束执行",
				Number:      task.Number + 1,
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("保存失败：%s", err.Error())
		}
	case CodeError:
		// 执行失败
		err := c.EditTask(c.gormClient.GetDb(), task.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行失败",
				Number:      task.Number + 1,
				RunId:       runId,
				UpdatedIp:   c.config.systemOutsideIp,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().AfterSeconds(task.Frequency).Time,
			}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("保存失败：%s", err.Error())
		}
	}

	if task.MaxNumber != 0 {
		if task.Number+1 >= task.MaxNumber {
			// 关闭执行
			err := c.EditTask(c.gormClient.GetDb(), task.Id).
				Select("status").
				Updates(jobs_gorm_model.Task{
					Status: TASK_TIMEOUT,
				}).Error
			if err != nil {
				c.zapLog.WithTraceId(ctx).Sugar().Errorf("保存失败：%s", err.Error())
			}
		}
	}
	return
}
