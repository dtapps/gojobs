package gojobs

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
)

// Ping 心跳
func (j *JobsGorm) Ping(ctx context.Context) {
	c := cron.New(cron.WithSeconds())
	_, _ = c.AddFunc(GetSeconds(2).Spec(), func() {
		result, err := j.redisClient.Set(ctx, j.config.cornKeyIp, j.config.outsideIp, 3).Result()
		if j.config.debug == true {
			log.Println("JOBS心跳", j.config.cornKeyIp, j.config.outsideIp, result, err)
		}
	})
	c.Start()
	defer c.Stop()
	select {}
}
