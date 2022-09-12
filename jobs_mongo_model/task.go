package jobs_mongo_model

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Task 任务
type Task struct {
	Id             primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`                          //【系统】记录编号
	Status         string             `json:"status,omitempty" bson:"status,omitempty"`                   //【系统】状态码
	Params         string             `json:"params,omitempty" bson:"params,omitempty"`                   //【系统】参数
	StatusDesc     string             `json:"status_desc,omitempty" bson:"status_desc,omitempty"`         //【系统】状态描述
	Frequency      int64              `json:"frequency,omitempty" bson:"frequency,omitempty"`             //【系统】频率(秒单位)
	Number         int64              `json:"number,omitempty" bson:"number,omitempty"`                   //【系统】当前次数
	MaxNumber      int64              `json:"max_number,omitempty" bson:"max_number,omitempty"`           //【系统】最大次数
	RunId          string             `json:"run_id,omitempty" bson:"run_id,omitempty"`                   //【系统】执行编号
	CustomId       string             `json:"custom_id,omitempty" bson:"custom_id,omitempty"`             //【系统】自定义编号
	CustomSequence int64              `json:"custom_sequence,omitempty" bson:"custom_sequence,omitempty"` //【系统】自定义顺序
	Type           string             `json:"type,omitempty" bson:"type,omitempty"`                       //【系统】类型
	TypeName       string             `json:"type_name,omitempty" bson:"type_name,omitempty"`             //【系统】类型名称
	SpecifyIp      string             `json:"specify_ip,omitempty" bson:"specify_ip,omitempty"`           //【系统】指定外网IP
	CreateRunInfo  struct {
		SystemHostName    string `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`       //【系统】主机名
		SystemInsideIp    string `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`       //【系统】内网ip
		SystemOs          string `json:"system_os,omitempty" bson:"system_os,omitempty"`                     //【系统】系统类型
		SystemArch        string `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                 //【系统】系统架构
		SystemCpuQuantity int    `json:"system_cpu_quantity,omitempty" bson:"system_cpu_quantity,omitempty"` //【系统】CPU核数
		GoVersion         string `json:"go_version,omitempty" bson:"go_version,omitempty"`                   //【程序】Go版本
		SdkVersion        string `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                 //【程序】Sdk版本
		RunTime           string `json:"run_time,omitempty" bson:"run_time,omitempty"`                       //【系统】运行时间
		RunIp             string `json:"run_ip,omitempty" bson:"run_ip,omitempty"`                           //【系统】外网ip
		RunResult         string `json:"run_result,omitempty" bson:"run_result,omitempty"`                   //【系统】结果
	} `json:"create_run_info,omitempty" bson:"create_run_info,omitempty"` //【系统】创建运行信息
	CurrentRunInfo struct {
		SystemHostName    string `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`       //【系统】主机名
		SystemInsideIp    string `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`       //【系统】内网ip
		SystemOs          string `json:"system_os,omitempty" bson:"system_os,omitempty"`                     //【系统】系统类型
		SystemArch        string `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                 //【系统】系统架构
		SystemCpuQuantity int    `json:"system_cpu_quantity,omitempty" bson:"system_cpu_quantity,omitempty"` //【系统】CPU核数
		GoVersion         string `json:"go_version,omitempty" bson:"go_version,omitempty"`                   //【程序】Go版本
		SdkVersion        string `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                 //【程序】Sdk版本
		RunTime           string `json:"run_time,omitempty" bson:"run_time,omitempty"`                       //【系统】运行时间
		RunIp             string `json:"run_ip,omitempty" bson:"run_ip,omitempty"`                           //【系统】外网ip
		RunResult         string `json:"run_result,omitempty" bson:"run_result,omitempty"`                   //【系统】结果
	} `json:"current_run_info,omitempty" bson:"current_run_info,omitempty"` //【系统】当前运行信息
	NextRunInfo struct {
		SystemHostName    string `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"`       //【系统】主机名
		SystemInsideIp    string `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"`       //【系统】内网ip
		SystemOs          string `json:"system_os,omitempty" bson:"system_os,omitempty"`                     //【系统】系统类型
		SystemArch        string `json:"system_arch,omitempty" bson:"system_arch,omitempty"`                 //【系统】系统架构
		SystemCpuQuantity int    `json:"system_cpu_quantity,omitempty" bson:"system_cpu_quantity,omitempty"` //【系统】CPU核数
		GoVersion         string `json:"go_version,omitempty" bson:"go_version,omitempty"`                   //【程序】Go版本
		SdkVersion        string `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`                 //【程序】Sdk版本
		RunTime           string `json:"run_time,omitempty" bson:"run_time,omitempty"`                       //【系统】运行时间
		RunIp             string `json:"run_ip,omitempty" bson:"run_ip,omitempty"`                           //【系统】外网ip
		RunResult         string `json:"run_result,omitempty" bson:"run_result,omitempty"`                   //【系统】结果
	} `json:"next_run_info,omitempty" bson:"next_run_info,omitempty"` //【系统】下一次运行信息
	CurrentTime primitive.DateTime `json:"current_time,omitempty" bson:"current_time,omitempty"` //【系统】创建时间
}

func (Task) TableName() string {
	return "task"
}
