package gojobs

import (
	"github.com/go-redis/redis/v9"
	"gorm.io/gorm"
)

// GetDb 获取数据库驱动
func (j *JobsGorm) GetDb() *gorm.DB {
	return j.gormClient.Db
}

// GetRedis 获取缓存数据库驱动
func (j *JobsGorm) GetRedis() *redis.Client {
	return j.redisClient.Db
}

// GetCurrentIp 获取当前ip
func (j *JobsGorm) GetCurrentIp() string {
	return j.config.outsideIp
}

// GetSubscribeAddress 获取订阅地址
func (j *JobsGorm) GetSubscribeAddress() string {
	return j.config.cornKeyPrefix + "_" + j.config.cornKeyCustom
}
