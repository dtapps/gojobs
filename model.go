package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gojobs/jobs_mongo_model"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 创建模型
func (c *Client) autoMigrateTask(ctx context.Context) {
	err := c.gormClient.GetDb().AutoMigrate(&jobs_gorm_model.Task{})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建模型：%s", err)
	}
}

// 创建模型
func (c *Client) autoMigrateTaskLog(ctx context.Context) {
	err := c.gormClient.GetDb().AutoMigrate(&jobs_gorm_model.TaskLog{})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建模型：%s", err)
	}
}

// GormTaskLogDelete 删除
func (c *Client) GormTaskLogDelete(ctx context.Context, hour int64) error {
	err := c.gormClient.GetDb().Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("删除失败：%s", err)
	}
	return err
}

// MongoTaskLogDelete 删除
func (c *Client) MongoTaskLogDelete(ctx context.Context, hour int64) (*mongo.DeleteResult, error) {
	filter := bson.D{{"log_time", bson.D{{"$lt", primitive.NewDateTimeFromTime(gotime.Current().BeforeHour(hour).Time)}}}}
	return c.mongoClient.Database(c.mongoConfig.databaseName).Collection(jobs_mongo_model.TaskLog{}.CollectionName()).DeleteMany(ctx, filter)
}

// 创建时间序列集合
func (c *Client) mongoCreateCollectionTaskLog(ctx context.Context) {
	err := c.mongoClient.Database(c.mongoConfig.databaseName).CreateCollection(ctx, jobs_mongo_model.TaskLog{}.CollectionName(), options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建时间序列集合：%s", err)
	}
}

// 创建索引
func (c *Client) mongoCreateIndexesTaskLog(ctx context.Context) {
	_, err := c.mongoClient.Database(c.mongoConfig.databaseName).Collection(jobs_mongo_model.TaskLog{}.CollectionName()).CreateManyIndexes(ctx, []mongo.IndexModel{{
		Keys: bson.D{{
			Key:   "log_time",
			Value: -1,
		}},
	}})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建索引：%s", err)
	}
}
