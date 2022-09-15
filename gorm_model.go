package gojobs

import (
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"gorm.io/gorm"
)

// TaskTakeId 查询单任务
func (c *Client) TaskTakeId(tx *gorm.DB, id uint) (result jobs_gorm_model.Task) {
	tx.Where("id = ?", id).Take(&result)
	return result
}

// TaskTake 查询单任务
func (c *Client) TaskTake(tx *gorm.DB, customId string) (result jobs_gorm_model.Task) {
	tx.Where("custom_id = ?", customId).Take(&result)
	return result
}

// 查询单任务
func (c *Client) taskTake(tx *gorm.DB, customId, status string) (result jobs_gorm_model.Task) {
	tx.Where("custom_id = ?", customId).Where("status = ?", status).Take(&result)
	return result
}

// TaskTakeIn 查询单任务 - 任务运行
func (c *Client) TaskTakeIn(tx *gorm.DB, customId string) jobs_gorm_model.Task {
	return c.taskTake(tx, customId, TASK_IN)
}

// TaskTakeSuccess 查询单任务 - 任务完成
func (c *Client) TaskTakeSuccess(tx *gorm.DB, customId string) jobs_gorm_model.Task {
	return c.taskTake(tx, customId, TASK_SUCCESS)
}

// TaskTakeError 查询单任务 - 任务异常
func (c *Client) TaskTakeError(tx *gorm.DB, customId string) jobs_gorm_model.Task {
	return c.taskTake(tx, customId, TASK_ERROR)
}

// TaskTakeTimeout 查询单任务 - 任务超时
func (c *Client) TaskTakeTimeout(tx *gorm.DB, customId string) jobs_gorm_model.Task {
	return c.taskTake(tx, customId, TASK_TIMEOUT)
}

// TaskTakeWait 查询单任务 - 任务等待
func (c *Client) TaskTakeWait(tx *gorm.DB, customId string) jobs_gorm_model.Task {
	return c.taskTake(tx, customId, TASK_WAIT)
}

// TaskTypeTake 查询单任务
func (c *Client) TaskTypeTake(tx *gorm.DB, customId, Type string) (result jobs_gorm_model.Task) {
	tx.Where("custom_id = ?", customId).Where("type = ?", Type).Take(&result)
	return result
}

// 查询单任务
func (c *Client) taskTypeTake(tx *gorm.DB, customId, Type, status string) (result jobs_gorm_model.Task) {
	tx.Where("custom_id = ?", customId).Where("type = ?", Type).Where("status = ?", status).Take(&result)
	return result
}

// TaskTypeTakeIn 查询单任务 - 任务运行
func (c *Client) TaskTypeTakeIn(tx *gorm.DB, customId, Type string) jobs_gorm_model.Task {
	return c.taskTypeTake(tx, customId, Type, TASK_IN)
}

// TaskTypeTakeSuccess 查询单任务 - 任务完成
func (c *Client) TaskTypeTakeSuccess(tx *gorm.DB, customId, Type string) jobs_gorm_model.Task {
	return c.taskTypeTake(tx, customId, Type, TASK_SUCCESS)
}

// TaskTypeTakeError 查询单任务 - 任务异常
func (c *Client) TaskTypeTakeError(tx *gorm.DB, customId, Type string) jobs_gorm_model.Task {
	return c.taskTypeTake(tx, customId, Type, TASK_ERROR)
}

// TaskTypeTakeTimeout 查询单任务 - 任务超时
func (c *Client) TaskTypeTakeTimeout(tx *gorm.DB, customId, Type string) jobs_gorm_model.Task {
	return c.taskTypeTake(tx, customId, Type, TASK_TIMEOUT)
}

// TaskTypeTakeWait 查询单任务 - 任务等待
func (c *Client) TaskTypeTakeWait(tx *gorm.DB, customId, Type string) jobs_gorm_model.Task {
	return c.taskTypeTake(tx, customId, Type, TASK_WAIT)
}

// TaskFindAll 查询多任务
func (c *Client) TaskFindAll(tx *gorm.DB, frequency int64) (results []jobs_gorm_model.Task) {
	tx.Where("frequency = ?", frequency).Order("id asc").Find(&results)
	return results
}

// 查询多任务
func (c *Client) taskFindAll(tx *gorm.DB, frequency int64, status string) (results []jobs_gorm_model.Task) {
	tx.Where("frequency = ?", frequency).Where("status = ?", status).Order("id asc").Find(&results)
	return results
}

// TaskFindAllIn 查询多任务 - 任务运行
func (c *Client) TaskFindAllIn(tx *gorm.DB, frequency int64) []jobs_gorm_model.Task {
	return c.taskFindAll(tx, frequency, TASK_IN)
}

// TaskFindAllSuccess 查询多任务 - 任务完成
func (c *Client) TaskFindAllSuccess(tx *gorm.DB, frequency int64) []jobs_gorm_model.Task {
	return c.taskFindAll(tx, frequency, TASK_SUCCESS)
}

// TaskFindAllError 查询多任务 - 任务异常
func (c *Client) TaskFindAllError(tx *gorm.DB, frequency int64) []jobs_gorm_model.Task {
	return c.taskFindAll(tx, frequency, TASK_ERROR)
}

// TaskFindAllTimeout 查询多任务 - 任务超时
func (c *Client) TaskFindAllTimeout(tx *gorm.DB, frequency int64) []jobs_gorm_model.Task {
	return c.taskFindAll(tx, frequency, TASK_TIMEOUT)
}

// TaskFindAllWait 查询多任务 - 任务等待
func (c *Client) TaskFindAllWait(tx *gorm.DB, frequency int64) []jobs_gorm_model.Task {
	return c.taskFindAll(tx, frequency, TASK_WAIT)
}

// StartTask 任务启动
func (c *Client) StartTask(tx *gorm.DB, id uint) error {
	return c.EditTask(tx, id).
		Select("status", "status_desc").
		Updates(jobs_gorm_model.Task{
			Status:     TASK_IN,
			StatusDesc: "启动任务",
		}).Error
}

// StartTaskCustom 任务启动自定义
func (c *Client) StartTaskCustom(tx *gorm.DB, customId string, customSequence int64) error {
	return tx.Model(&jobs_gorm_model.Task{}).
		Where("custom_id = ?", customId).
		Where("custom_sequence = ?", customSequence).
		Where("status = ?", TASK_WAIT).
		Select("status", "status_desc").
		Updates(jobs_gorm_model.Task{
			Status:     TASK_IN,
			StatusDesc: "启动任务",
		}).Error
}

// EditTask 任务修改
func (c *Client) EditTask(tx *gorm.DB, id uint) *gorm.DB {
	return tx.Model(&jobs_gorm_model.Task{}).Where("id = ?", id)
}

// UpdateFrequency 更新任务频率
func (c *Client) UpdateFrequency(tx *gorm.DB, id uint, frequency int64) *gorm.DB {
	return c.EditTask(tx, id).
		Select("frequency").
		Updates(jobs_gorm_model.Task{
			Frequency: frequency,
		})
}
