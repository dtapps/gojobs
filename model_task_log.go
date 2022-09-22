package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/golog"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TaskLog 任务日志模型
type TaskLog struct {
	LogId   primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"` //【记录】编号
	LogTime primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`    //【记录】时间
	Task    struct {
		Id         uint          `json:"id,omitempty" bson:"id,omitempty"`                   //【任务】编号
		RunId      string        `json:"run_id,omitempty" bson:"run_id,omitempty"`           //【任务】执行编号
		ResultCode int           `json:"result_code,omitempty" bson:"result_code,omitempty"` //【任务】执行状态码
		ResultDesc string        `json:"result_desc,omitempty" bson:"result_desc,omitempty"` //【任务】执行结果
		ResultTime dorm.BsonTime `json:"result_time,omitempty" bson:"result_time,omitempty"` //【任务】执行时间
	} `json:"task,omitempty" bson:"task,omitempty"` //【任务】信息
	System struct {
		HostName  string `json:"host_name,omitempty" bson:"host_name,omitempty"`   //【系统】主机名
		InsideIp  string `json:"inside_ip,omitempty" bson:"inside_ip,omitempty"`   //【系统】内网ip
		OutsideIp string `json:"outside_ip,omitempty" bson:"outside_ip,omitempty"` //【系统】外网ip
		Os        string `json:"os,omitempty" bson:"os,omitempty"`                 //【系统】系统类型
		Arch      string `json:"arch,omitempty" bson:"arch,omitempty"`             //【系统】系统架构
	} `json:"system,omitempty" bson:"system,omitempty"` //【系统】信息
	Version struct {
		Go  string `json:"go,omitempty" bson:"go,omitempty"`   //【程序】Go版本
		Sdk string `json:"sdk,omitempty" bson:"sdk,omitempty"` //【程序】Sdk版本
	} `json:"version,omitempty" bson:"version,omitempty"` //【程序】版本信息
}

func (TaskLog) CollectionName() string {
	return "task_log"
}

// 创建时间序列集合
func (TaskLog) createCollection(ctx context.Context, zapLog *golog.ZapLog, db *dorm.MongoClient, databaseName string) {
	err := db.Database(databaseName).CreateCollection(ctx, TaskLog{}.CollectionName(), options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
	if err != nil {
		zapLog.WithTraceId(ctx).Sugar().Errorf("创建时间序列集合：%s", err)
	}
}

// 创建索引
func (TaskLog) createIndexes(ctx context.Context, zapLog *golog.ZapLog, db *dorm.MongoClient, databaseName string) {
	_, err := db.Database(databaseName).Collection(TaskLog{}.CollectionName()).CreateManyIndexes(ctx, []mongo.IndexModel{{
		Keys: bson.D{{
			Key:   "log_time",
			Value: -1,
		}},
	}})
	if err != nil {
		zapLog.WithTraceId(ctx).Sugar().Errorf("创建索引：%s", err)
	}
}

// MongoTaskLogRecord 记录
func (c *Client) MongoTaskLogRecord(ctx context.Context, task jobs_gorm_model.Task, runId string, taskResultCode int, taskResultDesc string) {

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

	_, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(TaskLog{}.CollectionName()).InsertOne(ctx, taskLog)
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("记录失败：%s", err)
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("记录数据：%+v", taskLog)
	}

}

// MongoTaskLogDelete 删除
func (c *Client) MongoTaskLogDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Database(c.mongoConfig.databaseName).Collection(TaskLog{}.CollectionName()).DeleteMany(ctx, filter)
}
