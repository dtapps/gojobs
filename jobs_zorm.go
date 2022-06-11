package gojobs

import "gitee.com/chunanyong/zorm"

type JobsZorm struct {
	Db *zorm.DBDao
}

func NewJobsZorm(db *zorm.DBDao) *JobsZorm {
	var (
		jobsZorm = &JobsZorm{}
	)
	jobsZorm.Db = db
	return jobsZorm
}
