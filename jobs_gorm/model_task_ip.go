package jobs_gorm

import (
	"gorm.io/gorm"
	"log"
	"strings"
)

// TaskIp 任务Ip
type TaskIp struct {
	Id       int64  `gorm:"primaryKey" json:"id"` // 记录编号
	TaskType string `json:"task_type"`            // 任务编号
	Ips      string `json:"ips"`                  // 任务IP
}

func (m *TaskIp) TableName() string {
	return "task_ip"
}

func (jobsGorm *JobsGorm) TaskIpUpdate(tx *gorm.DB, taskType, ips string) *gorm.DB {
	var query TaskIp
	tx.Where("task_type = ?", taskType).Where("ips = ?", ips).Take(&query)
	if query.Id != 0 {
		return tx
	}
	updateStatus := tx.Create(&TaskIp{
		TaskType: taskType,
		Ips:      ips,
	})
	if updateStatus.RowsAffected == 0 {
		log.Println("任务更新失败：", updateStatus.Error)
	}
	return updateStatus
}

// TaskIpInit 实例任务ip
func (jobsGorm *JobsGorm) TaskIpInit(tx *gorm.DB, ips map[string]string) bool {
	if jobsGorm.outsideIp == "" || jobsGorm.outsideIp == "0.0.0.0" {
		return false
	}
	tx.Where("ips = ?", jobsGorm.outsideIp).Delete(&TaskIp{}) // 删除
	for k, v := range ips {
		if v == "" {
			jobsGorm.TaskIpUpdate(tx, k, jobsGorm.outsideIp)
		} else {
			find := strings.Contains(v, ",")
			if find == true {
				// 包含
				parts := strings.Split(v, ",")
				for _, vv := range parts {
					if vv == jobsGorm.outsideIp {
						jobsGorm.TaskIpUpdate(tx, k, jobsGorm.outsideIp)
					}
				}
			} else {
				// 不包含
				if v == jobsGorm.outsideIp {
					jobsGorm.TaskIpUpdate(tx, k, jobsGorm.outsideIp)
				}
			}
		}
	}
	return true
}
