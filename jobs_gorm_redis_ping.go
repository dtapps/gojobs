package gojobs

import (
	"context"
	"github.com/robfig/cron/v3"
	"time"
)

// Ping 心跳
func (j *JobsGorm) Ping(ctx context.Context) {
	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc(GetSeconds(2).Spec(), func() {
		result, err := j.redisClient.Set(ctx, j.config.cornKeyPrefix+"_"+j.config.cornKeyCustom, j.config.cornKeyCustom, 3*time.Second).Result()
		if j.config.logDebug == true {
			j.logClient.Infof(ctx, "[jobs.Ping] %s %s %v %s", j.config.cornKeyPrefix+"_"+j.config.cornKeyCustom, j.config.cornKeyCustom, result, err)
		}
	})
	c.Start()
	defer c.Stop()
	select {}
}
