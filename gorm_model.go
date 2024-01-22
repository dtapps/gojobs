package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
)

// TaskTakeId 编号查询任务
func (c *Client) TaskTakeId(ctx context.Context, tx *gorm.DB, id uint) (result gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("id = ?", id).
		Take(&result).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("编号查询任务：%v", err))
		}
	}
	return result
}

// TaskTake 自定义编号查询任务
func (c *Client) TaskTake(ctx context.Context, tx *gorm.DB, customId string) (result gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("custom_id = ?", customId).
		Take(&result).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("自定义编号查询任务：%v", err))
		}
	}
	return result
}

// 自定义编号加状态查询任务
func (c *Client) taskTake(ctx context.Context, tx *gorm.DB, customId, status string) (result gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("custom_id = ?", customId).
		Where("status = ?", status).
		Take(&result).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("自定义编号加状态查询任务：%v", err))
		}
	}
	return result
}

// TaskTakeIn 查询单任务 - 任务运行
func (c *Client) TaskTakeIn(ctx context.Context, tx *gorm.DB, customId string) gormModelTask {
	return c.taskTake(ctx, tx, customId, TASK_IN)
}

// TaskTakeSuccess 查询单任务 - 任务完成
func (c *Client) TaskTakeSuccess(ctx context.Context, tx *gorm.DB, customId string) gormModelTask {
	return c.taskTake(ctx, tx, customId, TASK_SUCCESS)
}

// TaskTakeError 查询单任务 - 任务异常
func (c *Client) TaskTakeError(ctx context.Context, tx *gorm.DB, customId string) gormModelTask {
	return c.taskTake(ctx, tx, customId, TASK_ERROR)
}

// TaskTakeTimeout 查询单任务 - 任务超时
func (c *Client) TaskTakeTimeout(ctx context.Context, tx *gorm.DB, customId string) gormModelTask {
	return c.taskTake(ctx, tx, customId, TASK_TIMEOUT)
}

// TaskTakeWait 查询单任务 - 任务等待
func (c *Client) TaskTakeWait(ctx context.Context, tx *gorm.DB, customId string) gormModelTask {
	return c.taskTake(ctx, tx, customId, TASK_WAIT)
}

// TaskTypeTake 查询单任务
func (c *Client) TaskTypeTake(ctx context.Context, tx *gorm.DB, customId, Type string) (result gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("custom_id = ?", customId).
		Where("type = ?", Type).
		Take(&result).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询单任务：%v", err))
		}
	}
	return result
}

// 查询单任务
func (c *Client) taskTypeTake(ctx context.Context, tx *gorm.DB, customId, Type, status string) (result gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("custom_id = ?", customId).Where("type = ?", Type).
		Where("status = ?", status).
		Take(&result).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询单任务：%v", err))
		}
	}
	return result
}

// TaskTypeTakeIn 查询单任务 - 任务运行
func (c *Client) TaskTypeTakeIn(ctx context.Context, tx *gorm.DB, customId, Type string) gormModelTask {
	return c.taskTypeTake(ctx, tx, customId, Type, TASK_IN)
}

// TaskTypeTakeSuccess 查询单任务 - 任务完成
func (c *Client) TaskTypeTakeSuccess(ctx context.Context, tx *gorm.DB, customId, Type string) gormModelTask {
	return c.taskTypeTake(ctx, tx, customId, Type, TASK_SUCCESS)
}

// TaskTypeTakeError 查询单任务 - 任务异常
func (c *Client) TaskTypeTakeError(ctx context.Context, tx *gorm.DB, customId, Type string) gormModelTask {
	return c.taskTypeTake(ctx, tx, customId, Type, TASK_ERROR)
}

// TaskTypeTakeTimeout 查询单任务 - 任务超时
func (c *Client) TaskTypeTakeTimeout(ctx context.Context, tx *gorm.DB, customId, Type string) gormModelTask {
	return c.taskTypeTake(ctx, tx, customId, Type, TASK_TIMEOUT)
}

// TaskTypeTakeWait 查询单任务 - 任务等待
func (c *Client) TaskTypeTakeWait(ctx context.Context, tx *gorm.DB, customId, Type string) gormModelTask {
	return c.taskTypeTake(ctx, tx, customId, Type, TASK_WAIT)
}

// TaskFindAll 查询多任务
func (c *Client) TaskFindAll(ctx context.Context, tx *gorm.DB, frequency int64) (results []gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("frequency = ?", frequency).
		Order("id asc").
		Find(&results).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询多任务：%v", err))
		}
	}
	return results
}

// TaskFindAllType 查询多任务
func (c *Client) TaskFindAllType(ctx context.Context, tx *gorm.DB, Type string, frequency int64) (results []gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("type = ?", Type).
		Where("frequency = ?", frequency).
		Order("id asc").
		Find(&results).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询多任务：%v", err))
		}
	}
	return results
}

// 查询多任务
func (c *Client) taskFindAll(ctx context.Context, tx *gorm.DB, frequency int64, status string) (results []gormModelTask) {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("frequency = ?", frequency).
		Where("status = ?", status).
		Order("id asc").
		Find(&results).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询多任务：%v", err))
		}
	}
	return results
}

// 查询多任务
func (c *Client) taskFindAllType(ctx context.Context, tx *gorm.DB, Type string, frequency int64, status string) (results []gormModelTask) {
	if frequency == 0 {
		err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
			Where("type = ?", Type).
			Where("status = ?", status).
			Order("id asc").
			Find(&results).Error
		if err != nil {
			if c.slog.status {
				c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询多任务：%v", err))
			}
		}
		return results
	}
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("type = ?", Type).
		Where("frequency = ?", frequency).
		Where("status = ?", status).
		Order("id asc").
		Find(&results).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("查询多任务：%v", err))
		}
	}
	return results
}

// TaskFindAllIn 查询多任务 - 任务运行
func (c *Client) TaskFindAllIn(ctx context.Context, tx *gorm.DB, frequency int64) []gormModelTask {
	return c.taskFindAll(ctx, tx, frequency, TASK_IN)
}

// TaskFindAllInType 查询多任务 - 任务运行
func (c *Client) TaskFindAllInType(ctx context.Context, tx *gorm.DB, Type string) []gormModelTask {
	return c.taskFindAllType(ctx, tx, Type, 0, TASK_IN)
}

// TaskFindAllSuccess 查询多任务 - 任务完成
func (c *Client) TaskFindAllSuccess(ctx context.Context, tx *gorm.DB, frequency int64) []gormModelTask {
	return c.taskFindAll(ctx, tx, frequency, TASK_SUCCESS)
}

// TaskFindAllSuccessType 查询多任务 - 任务完成
func (c *Client) TaskFindAllSuccessType(ctx context.Context, tx *gorm.DB, Type string) []gormModelTask {
	return c.taskFindAllType(ctx, tx, Type, 0, TASK_SUCCESS)
}

// TaskFindAllError 查询多任务 - 任务异常
func (c *Client) TaskFindAllError(ctx context.Context, tx *gorm.DB, frequency int64) []gormModelTask {
	return c.taskFindAll(ctx, tx, frequency, TASK_ERROR)
}

// TaskFindAllErrorType 查询多任务 - 任务异常
func (c *Client) TaskFindAllErrorType(ctx context.Context, tx *gorm.DB, Type string) []gormModelTask {
	return c.taskFindAllType(ctx, tx, Type, 0, TASK_ERROR)
}

// TaskFindAllTimeout 查询多任务 - 任务超时
func (c *Client) TaskFindAllTimeout(ctx context.Context, tx *gorm.DB, frequency int64) []gormModelTask {
	return c.taskFindAll(ctx, tx, frequency, TASK_TIMEOUT)
}

// TaskFindAllTimeoutType 查询多任务 - 任务超时
func (c *Client) TaskFindAllTimeoutType(ctx context.Context, tx *gorm.DB, Type string) []gormModelTask {
	return c.taskFindAllType(ctx, tx, Type, 0, TASK_TIMEOUT)
}

// TaskFindAllWait 查询多任务 - 任务等待
func (c *Client) TaskFindAllWait(ctx context.Context, tx *gorm.DB, frequency int64) []gormModelTask {
	return c.taskFindAll(ctx, tx, frequency, TASK_WAIT)
}

// TaskFindAllWaitType 查询多任务 - 任务等待
func (c *Client) TaskFindAllWaitType(ctx context.Context, tx *gorm.DB, Type string) []gormModelTask {
	return c.taskFindAllType(ctx, tx, Type, 0, TASK_WAIT)
}

// StartTask 任务启动
func (c *Client) StartTask(ctx context.Context, tx *gorm.DB, id uint) error {
	err := c.EditTask(ctx, tx, id).
		Select("status", "status_desc").
		Updates(gormModelTask{
			Status:     TASK_IN,
			StatusDesc: "启动任务",
		}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("任务启动失败：%v", err))
		}
	}
	return err
}

// StartTaskCustom 任务启动自定义
func (c *Client) StartTaskCustom(ctx context.Context, tx *gorm.DB, customId string, customSequence int64) error {
	err := tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("custom_id = ?", customId).
		Where("custom_sequence = ?", customSequence).
		Where("status = ?", TASK_WAIT).
		Select("status", "status_desc").
		Updates(gormModelTask{
			Status:     TASK_IN,
			StatusDesc: "启动任务",
		}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("任务启动自定义失败：%v", err))
		}
	}
	return err
}

// EditTask 任务修改
func (c *Client) EditTask(ctx context.Context, tx *gorm.DB, id uint) *gorm.DB {
	return tx.WithContext(ctx).Table(c.gormConfig.taskTableName).
		Where("id = ?", id)
}

// UpdateFrequency 更新任务频率
func (c *Client) UpdateFrequency(ctx context.Context, tx *gorm.DB, id uint, frequency int64) error {
	err := c.EditTask(ctx, tx, id).
		Select("frequency", "next_run_time").
		Updates(gormModelTask{
			Frequency:   frequency,
			NextRunTime: gotime.Current().AfterSeconds(frequency).Time,
		}).Error
	if err != nil {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("更新任务频率失败：%v", err))
		}
	}
	return err
}
