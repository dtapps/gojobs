package jobs_mongo_model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskLogRun 任务执行日志模型
type TaskLogRun struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                // 记录编号
	TaskId     uint               `json:"task_id,omitempty" bson:"task_id,omitempty"`       // 任务编号
	RunId      string             `json:"run_id,omitempty" bson:"run_id,omitempty"`         // 执行编号
	OutsideIp  string             `json:"outside_ip,omitempty" bson:"outside_ip,omitempty"` // 外网ip
	InsideIp   string             `json:"inside_ip,omitempty" bson:"inside_ip,omitempty"`   // 内网ip
	Os         string             `json:"os,omitempty" bson:"os,omitempty"`                 // 系统类型
	Arch       string             `json:"arch,omitempty" bson:"arch,omitempty"`             // 系统架构
	Gomaxprocs int                `json:"gomaxprocs,omitempty" bson:"gomaxprocs,omitempty"` // CPU核数
	GoVersion  string             `json:"go_version,omitempty" bson:"go_version,omitempty"` // GO版本
	MacAddrs   string             `json:"mac_addrs,omitempty" bson:"mac_addrs,omitempty"`   // Mac地址
	CreatedAt  primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"` // 创建时间
}

func (TaskLogRun) TableName() string {
	return "task_log_run"
}
