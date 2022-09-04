package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
)

// CheckManyTask 多任务检查
func (c *Client) CheckManyTask(ctx context.Context, tx *gorm.DB, vs []jobs_gorm_model.Task) {
	if len(vs) > 0 {
		for _, v := range vs {
			diffInSecondWithAbs := gotime.Current().DiffInSecondWithAbs(gotime.SetCurrent(v.UpdatedAt).Time)
			if diffInSecondWithAbs >= v.Frequency*3 {
				if c.config.debug == true {
					c.zapLog.WithTraceId(ctx).Sugar().Infof("[jobs.CheckManyTask]每隔%v秒任务：%v相差%v秒\n", v.Frequency, v.Id, diffInSecondWithAbs)
				}
				err := tx.Where("task_id = ?", v.Id).Where("run_id = ?", v.RunId).Delete(&jobs_gorm_model.TaskLogRun{}).Error
				if err != nil {
					c.zapLog.WithTraceId(ctx).Sugar().Errorf("删除失败：%s", err.Error())
				}
			}
		}
	}
}

// CheckSingleTask 单任务检查
func (c *Client) CheckSingleTask(ctx context.Context, tx *gorm.DB, v jobs_gorm_model.Task) {
	diffInSecondWithAbs := gotime.Current().DiffInSecondWithAbs(gotime.SetCurrent(v.UpdatedAt).Time)
	if diffInSecondWithAbs >= v.Frequency*3 {
		if c.config.debug == true {
			c.zapLog.WithTraceId(ctx).Sugar().Infof("[jobs.CheckManyTask]每隔%v秒任务：%v相差%v秒\n", v.Frequency, v.Id, diffInSecondWithAbs)
		}
		err := tx.Where("task_id = ?", v.Id).Where("run_id = ?", v.RunId).Delete(&jobs_gorm_model.TaskLogRun{}).Error
		if err != nil {
			c.zapLog.WithTraceId(ctx).Sugar().Errorf("删除失败：%s", err.Error())
		}
	}
}
