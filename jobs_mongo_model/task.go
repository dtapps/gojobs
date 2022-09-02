package jobs_mongo_model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task 任务
type Task struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                          // 记录编号
	Status         string             `json:"status,omitempty" bson:"status,omitempty"`                   // 状态码
	Params         string             `json:"params,omitempty" bson:"params,omitempty"`                   // 参数
	ParamsType     string             `json:"params_type,omitempty" bson:"params_type,omitempty"`         // 参数类型
	StatusDesc     string             `json:"status_desc,omitempty" bson:"status_desc,omitempty"`         // 状态描述
	Frequency      int64              `json:"frequency,omitempty" bson:"frequency,omitempty"`             // 频率(秒单位)
	Number         int64              `json:"number,omitempty" bson:"number,omitempty"`                   // 当前次数
	MaxNumber      int64              `json:"max_number,omitempty" bson:"max_number,omitempty"`           // 最大次数
	RunId          string             `json:"run_id,omitempty" bson:"run_id,omitempty"`                   // 执行编号
	CustomId       string             `json:"custom_id,omitempty" bson:"custom_id,omitempty"`             // 自定义编号
	CustomSequence int64              `json:"custom_sequence,omitempty" bson:"custom_sequence,omitempty"` // 自定义顺序
	Type           string             `json:"type,omitempty" bson:"type,omitempty"`                       // 类型
	TypeName       string             `json:"type_name,omitempty" bson:"type_name,omitempty"`             // 类型名称
	CreatedIp      string             `json:"created_ip,omitempty" bson:"created_ip,omitempty"`           // 创建外网IP
	SpecifyIp      string             `json:"specify_ip,omitempty" bson:"specify_ip,omitempty"`           // 指定外网IP
	UpdatedIp      string             `json:"updated_ip,omitempty" bson:"updated_ip,omitempty"`           // 更新外网IP
	Result         string             `json:"result,omitempty" bson:"result,omitempty"`                   // 结果
	NextRunTime    primitive.DateTime `json:"next_run_time,omitempty" bson:"next_run_time,omitempty"`     // 下次运行时间
	CreatedAt      primitive.DateTime `json:"created_at,omitempty" bson:"created_at,omitempty"`           // 创建时间
	UpdatedAt      primitive.DateTime `json:"updated_at,omitempty" bson:"updated_at,omitempty"`           // 更新时间
	DeletedAt      primitive.DateTime `json:"deleted_at,omitempty" bson:"deleted_at,omitempty"`           // 删除时间
}

func (Task) TableName() string {
	return "task"
}
