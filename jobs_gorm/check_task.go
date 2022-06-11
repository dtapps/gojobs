package jobs_gorm

import (
	"go.dtapp.net/gotime"
	"gorm.io/gorm"
	"log"
)

func (jobsGorm *JobsGorm) Check(tx *gorm.DB, vs []Task) {
	if jobsGorm.mainService > 0 && len(vs) > 0 {
		for _, v := range vs {
			diffInSecondWithAbs := gotime.Current().DiffInSecondWithAbs(gotime.SetCurrentParse(v.UpdatedAt).Time)
			if diffInSecondWithAbs >= v.Frequency*3 {
				log.Printf("每隔%v秒任务：%v相差%v秒\n", v.Frequency, v.Id, diffInSecondWithAbs)
				tx.Where("task_id = ?", v.Id).Where("run_id = ?", v.RunId).Delete(&TaskLogRun{}) // 删除
			}
		}
	}
}
