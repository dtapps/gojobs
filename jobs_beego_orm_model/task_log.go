package jobs_beego_orm_model

// TaskLog 任务日志模型
type TaskLog struct {
	Id         uint   `orm:"auto" json:"id"`                            // 记录编号
	TaskId     uint   `json:"task_id"`                                  // 任务编号
	StatusCode int    `json:"status_code"`                              // 状态码
	Desc       string `json:"desc"`                                     // 结果
	Version    int    `json:"version"`                                  // 版本
	CreatedAt  string `orm:"auto_now_add;type(text)" json:"created_at"` // 创建时间
}

func (m *TaskLog) TableName() string {
	return "task_log"
}
