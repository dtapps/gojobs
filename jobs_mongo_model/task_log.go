package jobs_mongo_model

import (
	"go.dtapp.net/dorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskLog 任务日志模型
type TaskLog struct {
	LogId           primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                          //【记录】编号
	LogTime         primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`                             //【记录】时间
	TaskId          uint               `json:"task_id,omitempty" bson:"task_id,omitempty"`                     //【任务】编号
	TaskRunId       string             `json:"task_run_id,omitempty" bson:"task_run_id,omitempty"`             //【任务】执行编号
	TaskResultCode  int                `json:"task_result_code,omitempty" bson:"task_result_code,omitempty"`   //【任务】执行状态码
	TaskResultDesc  string             `json:"task_result_desc,omitempty" bson:"task_result_desc,omitempty"`   //【任务】执行结果
	TaskResultTime  dorm.BsonTime      `json:"task_result_time,omitempty" bson:"task_result_time,omitempty"`   //【任务】执行时间
	SystemHostName  string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`   //【系统】主机名
	SystemInsideIp  string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`   //【系统】内网ip
	SystemOs        string             `json:"system_os,omitempty" bson:"system_os,omitempty"`                 //【系统】系统类型
	SystemArch      string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`             //【系统】系统架构
	GoVersion       string             `json:"go_version,omitempty" bson:"go_version,omitempty"`               //【系统】go版本
	SdkVersion      string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`             //【系统】sdk版本
	SystemOutsideIp string             `json:"system_outside_ip,omitempty" bson:"system_outside_ip,omitempty"` //【系统】外网ip
}

func (TaskLog) CollectionName() string {
	return "task_log"
}
