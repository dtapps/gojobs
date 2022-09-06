package jobs_mongo_model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskReceiveRecord 任务接收记录
type TaskReceiveRecord struct {
	Id         primitive.ObjectID        `json:"id,omitempty" bson:"_id,omitempty"`                  // 记录编号
	TaskInfo   TaskIssueRecordTaskInfo   `json:"task_info,omitempty" bson:"task_info,omitempty"`     // 任务信息
	SystemInfo TaskIssueRecordSystemInfo `json:"system_info,omitempty" bson:"system_info,omitempty"` // 系统信息
	RecordTime primitive.DateTime        `json:"record_time,omitempty" bson:"record_time,omitempty"` // 记录时间
}

func (TaskReceiveRecord) TableName() string {
	return "task_receive_record_"
}
