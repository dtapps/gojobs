package gojobs

import "gorm.io/gorm"

type JobsGorm struct {
	Db *gorm.DB
}

func NewJobsGorm(db *gorm.DB) *JobsGorm {
	var (
		jobsGorm = &JobsGorm{}
	)
	jobsGorm.Db = db
	return jobsGorm
}
