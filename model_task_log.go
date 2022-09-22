package gojobs

import (
	"go.dtapp.net/dorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// TaskLog 任务日志模型
type TaskLog struct {
	LogId   primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"` //【记录】编号
	LogTime primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`    //【记录】时间
	Task    struct {
		Id         uint          `json:"id,omitempty" bson:"id,omitempty"`                   //【任务】编号
		RunId      string        `json:"run_id,omitempty" bson:"run_id,omitempty"`           //【任务】执行编号
		ResultCode int           `json:"result_code,omitempty" bson:"result_code,omitempty"` //【任务】执行状态码
		ResultDesc string        `json:"result_desc,omitempty" bson:"result_desc,omitempty"` //【任务】执行结果
		ResultTime dorm.BsonTime `json:"result_time,omitempty" bson:"result_time,omitempty"` //【任务】执行时间
	} `json:"task,omitempty" bson:"task,omitempty"` //【任务】信息
	System struct {
		HostName  string `json:"host_name,omitempty" bson:"host_name,omitempty"`   //【系统】主机名
		InsideIp  string `json:"inside_ip,omitempty" bson:"inside_ip,omitempty"`   //【系统】内网ip
		OutsideIp string `json:"outside_ip,omitempty" bson:"outside_ip,omitempty"` //【系统】外网ip
		Os        string `json:"os,omitempty" bson:"os,omitempty"`                 //【系统】系统类型
		Arch      string `json:"arch,omitempty" bson:"arch,omitempty"`             //【系统】系统架构
	} `json:"system,omitempty" bson:"system,omitempty"` //【系统】信息
	Version struct {
		Go  string `json:"go,omitempty" bson:"go,omitempty"`   //【程序】Go版本
		Sdk string `json:"sdk,omitempty" bson:"sdk,omitempty"` //【程序】Sdk版本
	} `json:"version,omitempty" bson:"version,omitempty"` //【程序】版本信息
}

func (TaskLog) CollectionName() string {
	return "task_log"
}
