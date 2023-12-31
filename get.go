package gojobs

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// GetGormDb 获取数据库驱动
func (c *Client) GetGormDb() *gorm.DB {
	return c.gormClient
}

// GetRedisDb 获取缓存数据库驱动
func (c *Client) GetRedisDb() *redis.Client {
	return c.cache.redisClient
}

// GetCurrentIp 获取当前IP
func (c *Client) GetCurrentIp() string {
	return c.config.systemOutsideIp
}

// GetSubscribeAddress 获取订阅地址
func (c *Client) GetSubscribeAddress() string {
	return c.cache.cornKeyPrefix + "_" + c.cache.cornKeyCustom
}
