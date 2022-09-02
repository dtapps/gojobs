package jobs_xorm_model

// TaskIp 任务Ip
type TaskIp struct {
	Id       int64  `xorm:"pk autoincr" json:"id"`
	TaskType string `json:"task_type"` // 任务编号
	Ips      string `json:"ips"`       // 任务IP
}

func (TaskIp) TableName() string {
	return "task_ip"
}
