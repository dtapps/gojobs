package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gojobs/jobs_mongo_model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// 创建模型
func (c *Client) autoMigrateTask(ctx context.Context) {
	c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.gormClient.Db.AutoMigrate(&jobs_gorm_model.Task{}))
}

// 创建时间序列集合
func (c *Client) mongoCreateCollectionTask(ctx context.Context) {
	var commandResult bson.M
	commandErr := c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).RunCommand(ctx, bson.D{{
		"listCollections", 1,
	}}).Decode(&commandResult)
	if commandErr != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("检查时间序列集合：%s", commandErr)
	} else {
		c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).CreateCollection(ctx, jobs_mongo_model.Task{}.TableName(), options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("create_time"))))
	}
}

// 创建索引
func (c *Client) mongoCreateIndexesTask(ctx context.Context) {
	c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).Collection(jobs_mongo_model.Task{}.TableName()).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"status", 1},
	}}))
	c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).Collection(jobs_mongo_model.Task{}.TableName()).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"frequency", 1},
	}}))
	c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).Collection(jobs_mongo_model.Task{}.TableName()).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"custom_id", 1},
	}}))
	c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).Collection(jobs_mongo_model.Task{}.TableName()).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"type", 1},
	}}))
	c.zapLog.WithTraceId(ctx).Sugar().Info(c.db.mongoClient.Db.Database(c.db.mongoDatabaseName).Collection(jobs_mongo_model.Task{}.TableName()).Indexes().CreateOne(ctx, mongo.IndexModel{Keys: bson.D{
		{"specify_ip", 1},
	}}))
}
