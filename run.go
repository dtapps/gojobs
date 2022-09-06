package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gojobs/jobs_mongo_model"
	"go.dtapp.net/gostring"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Run 运行
func (c *Client) Run(ctx context.Context, info jobs_gorm_model.Task, status int, result string) {
	// 任务接收记录
	go func() {
		if c.cache.cornKeyCustom != "" {
			_, err := c.db.mongoClient.Database(c.db.mongoDatabaseName).
				Collection(jobs_mongo_model.TaskReceiveRecord{}.TableName() + c.cache.cornKeyCustom).
				InsertOne(&jobs_mongo_model.TaskReceiveRecord{
					Id: primitive.NewObjectID(),
					TaskInfo: jobs_mongo_model.TaskIssueRecordTaskInfo{
						Id:             info.Id,
						Status:         info.Status,
						Params:         info.Params,
						ParamsType:     info.ParamsType,
						StatusDesc:     info.StatusDesc,
						Frequency:      info.Frequency,
						Number:         info.Number,
						MaxNumber:      info.MaxNumber,
						RunId:          info.RunId,
						CustomId:       info.CustomId,
						CustomSequence: info.CustomSequence,
						Type:           info.Type,
						TypeName:       info.TypeName,
						CreatedIp:      info.CreatedIp,
						SpecifyIp:      info.SpecifyIp,
						UpdatedIp:      info.UpdatedIp,
						Result:         info.Result,
						NextRunTime:    dorm.BsonTime(info.NextRunTime),
						CreatedAt:      dorm.BsonTime(info.CreatedAt),
						UpdatedAt:      dorm.BsonTime(info.UpdatedAt),
					},
					SystemInfo: jobs_mongo_model.TaskIssueRecordSystemInfo{
						InsideIp:   c.config.insideIp,
						OutsideIp:  c.config.outsideIp,
						Os:         c.config.os,
						Arch:       c.config.arch,
						Gomaxprocs: c.config.maxProCs,
						GoVersion:  c.config.version,
						SdkVersion: c.config.runVersion,
					},
					RecordTime: primitive.NewDateTimeFromTime(gotime.Current().Time),
				})
			if err != nil {
				c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.jobs_mongo_model.TaskReceiveRecord]：%s", err.Error())
			}
		}
	}()
	// 请求函数记录
	err := c.db.gormClient.Db.Create(&jobs_gorm_model.TaskLog{
		TaskId:     info.Id,
		StatusCode: status,
		Desc:       result,
		Version:    c.config.runVersion,
	}).Error
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.Create]：%s", err.Error())
	}
	// 记录
	if c.db.mongoClient != nil && c.db.mongoClient.Db != nil {
		go func() {
			_, err = c.db.mongoClient.Database(c.db.mongoDatabaseName).
				Collection(jobs_mongo_model.TaskLog{}.TableName()).
				InsertOne(&jobs_mongo_model.TaskLog{
					Id:         primitive.NewObjectID(),
					TaskId:     info.Id,
					StatusCode: status,
					Desc:       result,
					Version:    c.config.runVersion,
					CreatedAt:  primitive.NewDateTimeFromTime(gotime.Current().Time),
				})
			if err != nil {
				c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.Run.jobs_mongo_model.TaskLog]：%s", err.Error())
			}
		}()
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
				UpdatedIp:   c.config.outsideIp,
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
				UpdatedIp:   c.config.outsideIp,
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
				UpdatedIp:   c.config.outsideIp,
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
	if c.db.mongoClient != nil && c.db.mongoClient.Db != nil {
		go func() {
			_, err := c.db.mongoClient.Database(c.db.mongoDatabaseName).
				Collection(jobs_mongo_model.TaskLogRun{}.TableName()).
				InsertOne(&jobs_mongo_model.TaskLogRun{
					Id:         primitive.NewObjectID(),
					TaskId:     id,
					RunId:      runId,
					InsideIp:   c.config.insideIp,
					OutsideIp:  c.config.outsideIp,
					Os:         c.config.os,
					Arch:       c.config.arch,
					Gomaxprocs: c.config.maxProCs,
					GoVersion:  c.config.version,
					SdkVersion: c.config.runVersion,
					MacAddrs:   c.config.macAddrS,
					CreatedAt:  primitive.NewDateTimeFromTime(gotime.Current().Time),
				})
			if err != nil {
				c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.RunAddLog.jobs_mongo_model.TaskLogRun]：%s", err.Error())
			}
		}()
	}
	return c.db.gormClient.Db.Create(&jobs_gorm_model.TaskLogRun{
		TaskId:     id,
		RunId:      runId,
		InsideIp:   c.config.insideIp,
		OutsideIp:  c.config.outsideIp,
		Os:         c.config.os,
		Arch:       c.config.arch,
		Gomaxprocs: c.config.maxProCs,
		GoVersion:  c.config.version,
		SdkVersion: c.config.runVersion,
		MacAddrs:   c.config.macAddrS,
	}).Error
}
