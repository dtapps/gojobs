package jobs_mongo_model

import "go.mongodb.org/mongo-driver/bson/primitive"

// TaskIp 任务Ip
type TaskIp struct {
	Id       primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`              // 记录编号
	TaskType string             `json:"task_type,omitempty" bson:"task_type,omitempty"` // 任务编号
	Ips      string             `json:"ips,omitempty" bson:"ips,omitempty"`             // 任务IP
}

func (TaskIp) TableName() string {
	return "task_ip"
}
