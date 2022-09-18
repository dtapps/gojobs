package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
)

// 创建模型
func (c *Client) autoMigrateTask(ctx context.Context) {
	err := c.gormClient.Db.AutoMigrate(&jobs_gorm_model.Task{})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建模型：%s", err)
	}
}

// 创建模型
func (c *Client) autoMigrateTaskLog(ctx context.Context) {
	err := c.gormClient.Db.AutoMigrate(&jobs_gorm_model.TaskLog{})
	if err != nil {
		c.zapLog.WithTraceId(ctx).Sugar().Errorf("创建模型：%s", err)
	}
}
