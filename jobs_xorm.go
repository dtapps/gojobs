package gojobs

import "xorm.io/xorm"

type JobsXorm struct {
	Db *xorm.Engine
}

func newJobsXorm(db *xorm.Engine) *JobsXorm {
	var (
		jobsXorm = &JobsXorm{}
	)
	jobsXorm.Db = db
	return jobsXorm
}
