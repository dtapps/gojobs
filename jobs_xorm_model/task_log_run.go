package jobs_xorm_model

// TaskLogRun 任务执行日志模型
type TaskLogRun struct {
	Id         uint   `xorm:"pk autoincr" json:"id"`     // 记录编号
	TaskId     uint   `json:"task_id"`                   // 任务编号
	RunId      string `json:"run_id"`                    // 执行编号
	OutsideIp  string `json:"outside_ip"`                // 外网ip
	InsideIp   string `json:"inside_ip"`                 // 内网ip
	Os         string `json:"os"`                        // 系统类型
	Arch       string `json:"arch"`                      // 系统架构
	Gomaxprocs int    `json:"gomaxprocs"`                // CPU核数
	GoVersion  string `json:"go_version"`                // GO版本
	SdkVersion string `json:"sdk_version"`               // SDK版本
	MacAddrs   string `json:"mac_addrs"`                 // Mac地址
	CreatedAt  string `xorm:"created" json:"created_at"` // 创建时间
}

func (TaskLogRun) TableName() string {
	return "task_log_run"
}
