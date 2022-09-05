package jobs_mongo_model

import (
	"go.dtapp.net/dorm"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TaskIssueRecordTaskInfo struct {
	Id             uint          `json:"id,omitempty" bson:"id,omitempty"`                           // 记录编号
	Status         string        `json:"status,omitempty" bson:"status,omitempty"`                   // 状态码
	Params         string        `json:"params,omitempty" bson:"params,omitempty"`                   // 参数
	ParamsType     string        `json:"params_type,omitempty" bson:"params_type,omitempty"`         // 参数类型
	StatusDesc     string        `json:"status_desc,omitempty" bson:"status_desc,omitempty"`         // 状态描述
	Frequency      int64         `json:"frequency,omitempty" bson:"frequency,omitempty"`             // 频率(秒单位)
	Number         int64         `json:"number,omitempty" bson:"number,omitempty"`                   // 当前次数
	MaxNumber      int64         `json:"max_number,omitempty" bson:"max_number,omitempty"`           // 最大次数
	RunId          string        `json:"run_id,omitempty" bson:"run_id,omitempty"`                   // 执行编号
	CustomId       string        `json:"custom_id,omitempty" bson:"custom_id,omitempty"`             // 自定义编号
	CustomSequence int64         `json:"custom_sequence,omitempty" bson:"custom_sequence,omitempty"` // 自定义顺序
	Type           string        `json:"type,omitempty" bson:"type,omitempty"`                       // 类型
	TypeName       string        `json:"type_name,omitempty" bson:"type_name,omitempty"`             // 类型名称
	CreatedIp      string        `json:"created_ip,omitempty" bson:"created_ip,omitempty"`           // 创建外网IP
	SpecifyIp      string        `json:"specify_ip,omitempty" bson:"specify_ip,omitempty"`           // 指定外网IP
	UpdatedIp      string        `json:"updated_ip,omitempty" bson:"updated_ip,omitempty"`           // 更新外网IP
	Result         string        `json:"result,omitempty" bson:"result,omitempty"`                   // 结果
	NextRunTime    dorm.BsonTime `json:"next_run_time,omitempty" bson:"next_run_time,omitempty"`     // 下次运行时间
	CreatedAt      dorm.BsonTime `json:"created_at,omitempty" bson:"created_at,omitempty"`           // 创建时间
	UpdatedAt      dorm.BsonTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`           // 更新时间
}

type TaskIssueRecordSystemInfo struct {
	OutsideIp  string `json:"outside_ip,omitempty" bson:"outside_ip,omitempty"`   // 外网ip
	InsideIp   string `json:"inside_ip,omitempty" bson:"inside_ip,omitempty"`     // 内网ip
	Os         string `json:"os,omitempty" bson:"os,omitempty"`                   // 系统类型
	Arch       string `json:"arch,omitempty" bson:"arch,omitempty"`               // 系统架构
	Gomaxprocs int    `json:"gomaxprocs,omitempty" bson:"gomaxprocs,omitempty"`   // CPU核数
	GoVersion  string `json:"go_version,omitempty" bson:"go_version,omitempty"`   // GO版本
	SdkVersion string `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"` // SDK版本
}

// TaskIssueRecord 任务发布记录
type TaskIssueRecord struct {
	Id            primitive.ObjectID        `json:"id,omitempty" bson:"_id,omitempty"`                        // 记录编号
	TaskInfo      TaskIssueRecordTaskInfo   `json:"task_info,omitempty" bson:"task_info,omitempty"`           // 任务信息
	SystemInfo    TaskIssueRecordSystemInfo `json:"system_info,omitempty" bson:"system_info,omitempty"`       // 系统信息
	RecordAddress string                    `json:"record_address,omitempty" bson:"record_address,omitempty"` // 接收地址
	RecordTime    dorm.BsonTime             `json:"record_time,omitempty" bson:"record_time,omitempty"`       // 记录时间
}

func (TaskIssueRecord) TableName() string {
	return "task_issue_record"
}
