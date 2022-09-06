package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gojobs/jobs_mongo_model"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// PublishLog 发布记录
func (c *Client) PublishLog(ctx context.Context, info jobs_gorm_model.Task, recordAddress string) {
	_, err := c.db.mongoClient.Database(c.db.mongoDatabaseName).
		Collection(jobs_mongo_model.TaskIssueRecord{}.TableName()).
		InsertOne(&jobs_mongo_model.TaskIssueRecord{
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
			RecordAddress: recordAddress,
			RecordTime:    primitive.NewDateTimeFromTime(gotime.Current().Time),
		})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("[gojobs.RunAddLog.jobs_mongo_model.TaskIssueRecord]：%s", err.Error())
	}
}
