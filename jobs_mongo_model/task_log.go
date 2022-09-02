package jobs_mongo_model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskLog 任务日志模型
type TaskLog struct {
	Id         primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                  // 记录编号
	TaskId     uint               `json:"task_id,omitempty" bson:"task_id,omitempty"`         // 任务编号
	StatusCode int                `json:"status_code,omitempty" bson:"status_code,omitempty"` // 状态码
	Desc       string             `json:"desc,omitempty" bson:"desc,omitempty"`               // 结果
	Version    string             `json:"version,omitempty" bson:"version,omitempty"`         // 版本
	CreatedAt  primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`   // 创建时间
}

func (TaskLog) TableName() string {
	return "task_log"
}
