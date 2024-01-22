package gojobs

import "time"

// 任务日志
type gormModelTaskLog struct {
	LogID           uint      `gorm:"primaryKey;comment:【日志】编号" json:"log_id"`                            // 【日志】编号
	LogTime         time.Time `gorm:"autoCreateTime;index;comment:【日志】时间" json:"log_time"`                // 【日志】时间
	TaskID          uint      `gorm:"index;comment:【任务】编号" json:"task_id"`                                // 【任务】编号
	TaskRunID       string    `gorm:"comment:【任务】执行编号" json:"task_run_id"`                                // 【任务】执行编号
	TaskResultCode  int       `gorm:"index;comment:【任务】执行状态码" json:"task_result_code"`                    // 【任务】执行状态码
	TaskResultDesc  string    `gorm:"comment:【任务】执行结果" json:"task_result_desc"`                           // 【任务】执行结果
	SystemHostName  string    `gorm:"comment:【系统】主机名" json:"system_host_name,omitempty"`                  //【系统】主机名
	SystemInsideIP  string    `gorm:"default:0.0.0.0;comment:【系统】内网IP" json:"system_inside_ip,omitempty"` //【系统】内网IP
	SystemOutsideIP string    `gorm:"default:0.0.0.0;comment:【系统】外网IP" json:"system_outside_ip"`          //【系统】外网IP
	SystemOs        string    `gorm:"comment:【系统】类型" json:"system_os,omitempty"`                          //【系统】类型
	SystemArch      string    `gorm:"comment:【系统】架构" json:"system_arch,omitempty"`                        //【系统】架构
	SystemUpTime    uint64    `gorm:"comment:【系统】运行时间" json:"system_up_time,omitempty"`                   //【系统】运行时间
	SystemBootTime  uint64    `gorm:"comment:【系统】开机时间" json:"system_boot_time,omitempty"`                 //【系统】开机时间
	GoVersion       string    `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`                       //【程序】Go版本
	SdkVersion      string    `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`                     //【程序】Sdk版本
	SystemVersion   string    `gorm:"comment:【程序】System版本" json:"system_version,omitempty"`               //【程序】System版本
	CpuCores        int       `gorm:"comment:【CPU】核数" json:"cpu_cores,omitempty"`                         //【CPU】核数
	CpuModelName    string    `gorm:"comment:【CPU】型号名称" json:"cpu_model_name,omitempty"`                  //【CPU】型号名称
	CpuMhz          float64   `gorm:"comment:【CPU】兆赫" json:"cpu_mhz,omitempty"`                           //【CPU】兆赫
}
