package jobs_gorm

import (
	"go.dtapp.net/gojobs/jobs_common"
	"gorm.io/gorm"
)

// Task 任务
type Task struct {
	Id             uint           `gorm:"primaryKey;comment:记录编号" json:"id"`              // 记录编号
	Status         string         `gorm:"index;comment:状态码" json:"status"`                // 状态码
	Params         string         `gorm:"comment:参数" json:"params"`                       // 参数
	ParamsType     string         `gorm:"comment:参数类型" json:"params_type"`                // 参数类型
	StatusDesc     string         `gorm:"comment:状态描述" json:"status_desc"`                // 状态描述
	Frequency      int64          `gorm:"index;comment:频率(秒单位)" json:"frequency"`         // 频率(秒单位)
	Number         int64          `gorm:"comment:当前次数" json:"number"`                     // 当前次数
	MaxNumber      int64          `gorm:"comment:最大次数" json:"max_number"`                 // 最大次数
	RunId          string         `gorm:"comment:执行编号" json:"run_id"`                     // 执行编号
	CustomId       string         `gorm:"index;comment:自定义编号" json:"custom_id"`           // 自定义编号
	CustomSequence int64          `gorm:"comment:自定义顺序" json:"custom_sequence"`           // 自定义顺序
	Type           string         `gorm:"index;comment:类型" json:"type"`                   // 类型
	CreatedIp      string         `gorm:"comment:创建外网IP" json:"created_ip"`               // 创建外网IP
	SpecifyIp      string         `gorm:"comment:指定外网IP" json:"specify_ip"`               // 指定外网IP
	UpdatedIp      string         `gorm:"comment:更新外网IP" json:"updated_ip"`               // 更新外网IP
	Result         string         `gorm:"comment:结果" json:"result"`                       // 结果
	CreatedAt      string         `gorm:"type:text;comment:创建时间" json:"created_at"`       // 创建时间
	UpdatedAt      string         `gorm:"type:text;comment:更新时间" json:"updated_at"`       // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"type:text;index;comment:删除时间" json:"deleted_at"` // 删除时间
}

func (m *Task) TableName() string {
	return "task"
}

// TaskTake 查询单任务
func (jobsGorm *JobsGorm) TaskTake(tx *gorm.DB, customId string) (result Task) {
	tx.Where("custom_id = ?", customId).Take(&result)
	return result
}

// 查询单任务
func (jobsGorm *JobsGorm) taskTake(tx *gorm.DB, customId, status string) (result Task) {
	tx.Where("custom_id = ?", customId).Where("status = ?", status).Take(&result)
	return result
}

// TaskTakeIn 查询单任务 - 任务运行
func (jobsGorm *JobsGorm) TaskTakeIn(tx *gorm.DB, customId string) Task {
	return jobsGorm.taskTake(tx, customId, jobs_common.TASK_IN)
}

// TaskTakeSuccess 查询单任务 - 任务完成
func (jobsGorm *JobsGorm) TaskTakeSuccess(tx *gorm.DB, customId string) Task {
	return jobsGorm.taskTake(tx, customId, jobs_common.TASK_SUCCESS)
}

// TaskTakeError 查询单任务 - 任务异常
func (jobsGorm *JobsGorm) TaskTakeError(tx *gorm.DB, customId string) Task {
	return jobsGorm.taskTake(tx, customId, jobs_common.TASK_ERROR)
}

// TaskTakeTimeout 查询单任务 - 任务超时
func (jobsGorm *JobsGorm) TaskTakeTimeout(tx *gorm.DB, customId string) Task {
	return jobsGorm.taskTake(tx, customId, jobs_common.TASK_TIMEOUT)
}

// TaskTakeWait 查询单任务 - 任务等待
func (jobsGorm *JobsGorm) TaskTakeWait(tx *gorm.DB, customId string) Task {
	return jobsGorm.taskTake(tx, customId, jobs_common.TASK_WAIT)
}

// TaskTypeTake 查询单任务
func (jobsGorm *JobsGorm) TaskTypeTake(tx *gorm.DB, customId, Type string) (result Task) {
	tx.Where("custom_id = ?", customId).Where("type = ?", Type).Take(&result)
	return result
}

// 查询单任务
func (jobsGorm *JobsGorm) taskTypeTake(tx *gorm.DB, customId, Type, status string) (result Task) {
	tx.Where("custom_id = ?", customId).Where("type = ?", Type).Where("status = ?", status).Take(&result)
	return result
}

// TaskTypeTakeIn 查询单任务 - 任务运行
func (jobsGorm *JobsGorm) TaskTypeTakeIn(tx *gorm.DB, customId, Type string) Task {
	return jobsGorm.taskTypeTake(tx, customId, Type, jobs_common.TASK_IN)
}

// TaskTypeTakeSuccess 查询单任务 - 任务完成
func (jobsGorm *JobsGorm) TaskTypeTakeSuccess(tx *gorm.DB, customId, Type string) Task {
	return jobsGorm.taskTypeTake(tx, customId, Type, jobs_common.TASK_SUCCESS)
}

// TaskTypeTakeError 查询单任务 - 任务异常
func (jobsGorm *JobsGorm) TaskTypeTakeError(tx *gorm.DB, customId, Type string) Task {
	return jobsGorm.taskTypeTake(tx, customId, Type, jobs_common.TASK_ERROR)
}

// TaskTypeTakeTimeout 查询单任务 - 任务超时
func (jobsGorm *JobsGorm) TaskTypeTakeTimeout(tx *gorm.DB, customId, Type string) Task {
	return jobsGorm.taskTypeTake(tx, customId, Type, jobs_common.TASK_TIMEOUT)
}

// TaskTypeTakeWait 查询单任务 - 任务等待
func (jobsGorm *JobsGorm) TaskTypeTakeWait(tx *gorm.DB, customId, Type string) Task {
	return jobsGorm.taskTypeTake(tx, customId, Type, jobs_common.TASK_WAIT)
}

// TaskFindAll 查询多任务
func (jobsGorm *JobsGorm) TaskFindAll(tx *gorm.DB, frequency int64) (results []Task) {
	tx.Where("frequency = ?", frequency).Order("id asc").Find(&results)
	return results
}

// 查询多任务
func (jobsGorm *JobsGorm) taskFindAll(tx *gorm.DB, frequency int64, status string) (results []Task) {
	tx.Where("frequency = ?", frequency).Where("status = ?", status).Order("id asc").Find(&results)
	return results
}

// TaskFindAllIn 查询多任务 - 任务运行
func (jobsGorm *JobsGorm) TaskFindAllIn(tx *gorm.DB, frequency int64) []Task {
	return jobsGorm.taskFindAll(tx, frequency, jobs_common.TASK_IN)
}

// TaskFindAllSuccess 查询多任务 - 任务完成
func (jobsGorm *JobsGorm) TaskFindAllSuccess(tx *gorm.DB, frequency int64) []Task {
	return jobsGorm.taskFindAll(tx, frequency, jobs_common.TASK_SUCCESS)
}

// TaskFindAllError 查询多任务 - 任务异常
func (jobsGorm *JobsGorm) TaskFindAllError(tx *gorm.DB, frequency int64) []Task {
	return jobsGorm.taskFindAll(tx, frequency, jobs_common.TASK_ERROR)
}

// TaskFindAllTimeout 查询多任务 - 任务超时
func (jobsGorm *JobsGorm) TaskFindAllTimeout(tx *gorm.DB, frequency int64) []Task {
	return jobsGorm.taskFindAll(tx, frequency, jobs_common.TASK_TIMEOUT)
}

// TaskFindAllWait 查询多任务 - 任务等待
func (jobsGorm *JobsGorm) TaskFindAllWait(tx *gorm.DB, frequency int64) []Task {
	return jobsGorm.taskFindAll(tx, frequency, jobs_common.TASK_WAIT)
}
