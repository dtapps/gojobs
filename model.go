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
	indexes, err := c.db.mongoClient.Database(c.db.mongoDatabaseName).Collection(jobs_mongo_model.Task{}.TableName()).CreateManyIndexes(ctx, []mongo.IndexModel{
		{
			Keys: bson.D{{
				Key:   "status",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "frequency",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "custom_id",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "type",
				Value: 1,
			}},
		}, {
			Keys: bson.D{{
				Key:   "specify_ip",
				Value: 1,
			}},
		},
	})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建索引：%s", err)
	}
	c.zapLog.WithTraceId(ctx).Sugar().Infof("创建索引：%s", indexes)
}
