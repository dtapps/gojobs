package jobs_gorm

import (
	"go.dtapp.net/gojobs/jobs_common"
	"gorm.io/gorm"
)

// Task 任务
type Task struct {
	Id             uint           `gorm:"primaryKey" json:"id"`              // 记录编号
	Status         string         `json:"status"`                            // 状态码
	Params         string         `json:"params"`                            // 参数
	ParamsType     string         `json:"params_type"`                       // 参数类型
	StatusDesc     string         `json:"status_desc"`                       // 状态描述
	Frequency      int64          `json:"frequency"`                         // 频率（秒单位）
	Number         int64          `json:"number"`                            // 当前次数
	MaxNumber      int64          `json:"max_number"`                        // 最大次数
	RunId          string         `json:"run_id"`                            // 执行编号
	CustomId       string         `json:"custom_id"`                         // 自定义编号
	CustomSequence int64          `json:"custom_sequence"`                   // 自定义顺序
	Type           string         `json:"type"`                              // 类型
	CreatedIp      string         `json:"created_ip"`                        // 创建外网IP
	SpecifyIp      string         `json:"specify_ip"`                        // 指定外网IP
	UpdatedIp      string         `json:"updated_ip"`                        // 更新外网IP
	Result         string         `json:"result"`                            // 结果
	CreatedAt      string         `gorm:"type:text" json:"created_at"`       // 创建时间
	UpdatedAt      string         `gorm:"type:text" json:"updated_at"`       // 更新时间
	DeletedAt      gorm.DeletedAt `gorm:"type:text;index" json:"deleted_at"` // 删除时间
}

func (m *Task) TableName() string {
	return "task"
}

// TaskTake 查询任务
func (jobsGorm *JobsGorm) TaskTake(tx *gorm.DB, customId string) (result Task) {
	tx.Where("custom_id = ?", customId).Where("status = ?", jobs_common.TASK_IN).Take(&result)
	return result
}

// TaskCustomIdTake 查询任务
func (jobsGorm *JobsGorm) TaskCustomIdTake(tx *gorm.DB, Type, customId string) (result Task) {
	tx.Where("type = ?", Type).Where("custom_id = ?", customId).Take(&result)
	return result
}

// TaskCustomIdTakeStatus 查询任务
func (jobsGorm *JobsGorm) TaskCustomIdTakeStatus(tx *gorm.DB, Type, customId, status string) (result Task) {
	tx.Where("type = ?", Type).Where("custom_id = ?", customId).Where("status = ?", status).Take(&result)
	return result
}

// TaskFind 查询任务
func (jobsGorm *JobsGorm) TaskFind(tx *gorm.DB, frequency int64) (results []Task) {
	tx.Table("task").Select("task.*").Where("task.frequency = ?", frequency).Where("task.status = ?", jobs_common.TASK_IN).Where("task_ip.ips = ?", jobsGorm.outsideIp).Order("task.id asc").Joins("left join task_ip on task_ip.task_type = task.type").Find(&results)
	return jobsGorm.taskFindCheck(results)
}

// TaskFindAll 查询任务
func (jobsGorm *JobsGorm) TaskFindAll(tx *gorm.DB, frequency int64) (results []Task) {
	tx.Where("frequency = ?", frequency).Where("status = ?", jobs_common.TASK_IN).Order("id asc").Find(&results)
	return results
}

// 检查任务
func (jobsGorm *JobsGorm) taskFindCheck(lists []Task) (results []Task) {
	for _, v := range lists {
		if v.SpecifyIp == "" {
			results = append(results, v)
		} else {
			if jobsGorm.outsideIp == v.SpecifyIp {
				results = append(results, v)
			}
		}
	}
	return results
}
