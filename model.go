package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
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

// TaskLogDelete 删除
func (c *Client) TaskLogDelete(ctx context.Context, hour int64) error {
	return c.gormClient.Db.Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).Delete(&jobs_gorm_model.TaskLog{}).Error
}
