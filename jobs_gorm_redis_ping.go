package gojobs

import (
	"context"
	"log"
)

// Ping 心跳
func (j *JobsGorm) Ping(ctx context.Context) error {
	result, err := j.redisClient.Set(ctx, j.config.cornKeyIp, j.config.outsideIp, 5).Result()

	if j.config.debug == true {
		log.Println("gojobs.Ping", j.config.cornKeyIp, j.config.outsideIp, result, err)
	}

	return err
}
