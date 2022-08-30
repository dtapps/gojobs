package gojobs

import (
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gostring"
	"log"
)

// Run 运行
func (j *JobsGorm) Run(info jobs_gorm_model.Task, status int, desc string) {
	// 请求函数记录
	err := j.gormClient.Db.Create(&jobs_gorm_model.TaskLog{
		TaskId:     info.Id,
		StatusCode: status,
		Desc:       desc,
		Version:    j.config.runVersion,
	}).Error
	if err != nil {
		log.Println("[gojobs.Run.Create]", err.Error())
	}
	if status == 0 {
		err = j.EditTask(j.gormClient.Db, info.Id).
			Select("run_id").
			Updates(jobs_gorm_model.Task{
				RunId: gostring.GetUuId(),
			}).Error
		if err != nil {
			log.Println("[gojobs.Run.EditTask]", err.Error())
		}
		return
	}
	// 任务
	if status == CodeSuccess {
		// 执行成功
		err = j.EditTask(j.gormClient.Db, info.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result").
			Updates(jobs_gorm_model.Task{
				StatusDesc: "执行成功",
				Number:     info.Number + 1,
				RunId:      gostring.GetUuId(),
				UpdatedIp:  j.config.outsideIp,
				Result:     desc,
			}).Error
		if err != nil {
			log.Println("[gojobs.Run.EditTask]", err.Error())
		}
	}
	if status == CodeEnd {
		// 执行成功、提前结束
		err = j.EditTask(j.gormClient.Db, info.Id).
			Select("status", "status_desc", "number", "updated_ip", "result").
			Updates(jobs_gorm_model.Task{
				Status:     TASK_SUCCESS,
				StatusDesc: "结束执行",
				Number:     info.Number + 1,
				UpdatedIp:  j.config.outsideIp,
				Result:     desc,
			}).Error
		if err != nil {
			log.Println("[gojobs.Run.EditTask]", err.Error())
		}
	}
	if status == CodeError {
		// 执行失败
		err = j.EditTask(j.gormClient.Db, info.Id).
			Select("status_desc", "number", "run_id", "updated_ip", "result").
			Updates(jobs_gorm_model.Task{
				StatusDesc: "执行失败",
				Number:     info.Number + 1,
				RunId:      gostring.GetUuId(),
				UpdatedIp:  j.config.outsideIp,
				Result:     desc,
			}).Error
		if err != nil {
			log.Println("[gojobs.Run.EditTask]", err.Error())
		}
	}
	if info.MaxNumber != 0 {
		if info.Number+1 >= info.MaxNumber {
			// 关闭执行
			err = j.EditTask(j.gormClient.Db, info.Id).
				Select("status").
				Updates(jobs_gorm_model.Task{
					Status: TASK_TIMEOUT,
				}).Error
			if err != nil {
				log.Println("[gojobs.Run.EditTask]", err.Error())
			}
		}
	}
}

// RunAddLog 任务执行日志
func (j *JobsGorm) RunAddLog(id uint, runId string) error {
	return j.gormClient.Db.Create(&jobs_gorm_model.TaskLogRun{
		TaskId:     id,
		RunId:      runId,
		InsideIp:   j.config.insideIp,
		OutsideIp:  j.config.outsideIp,
		Os:         j.config.os,
		Arch:       j.config.arch,
		Gomaxprocs: j.config.maxProCs,
		GoVersion:  j.config.version,
		MacAddrs:   j.config.macAddrS,
	}).Error
}
