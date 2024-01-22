package gojobs

import (
	"context"
	"errors"
	"github.com/redis/go-redis/v9"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	"runtime"
)

type systemResult struct {
	SystemHostname      string  // 主机名
	SystemOs            string  // 系统类型
	SystemVersion       string  // 系统版本
	SystemKernel        string  // 系统内核
	SystemKernelVersion string  // 系统内核版本
	SystemUpTime        uint64  // 系统运行时间
	SystemBootTime      uint64  // 系统开机时间
	CpuCores            int     // CPU核数
	CpuModelName        string  // CPU型号名称
	CpuMhz              float64 // CPU兆赫
}

// 获取系统信息
func getSystem() (result systemResult) {

	hInfo, _ := host.Info()

	result.SystemHostname = hInfo.Hostname
	result.SystemOs = hInfo.OS
	result.SystemVersion = hInfo.PlatformVersion
	result.SystemKernel = hInfo.KernelArch
	result.SystemKernelVersion = hInfo.KernelVersion
	result.SystemUpTime = hInfo.Uptime
	if hInfo.BootTime != 0 {
		result.SystemBootTime = hInfo.BootTime
	}

	hCpu, _ := cpu.Times(true)

	result.CpuCores = len(hCpu)

	cInfo, _ := cpu.Info()

	if len(cInfo) > 0 {
		result.CpuModelName = cInfo[0].ModelName
		result.CpuMhz = cInfo[0].Mhz
	}

	return result
}

// 设置配置信息
func (c *Client) setConfig(ctx context.Context, systemOutsideIP string) {

	info := getSystem()

	c.config.systemHostname = info.SystemHostname
	c.config.systemOs = info.SystemOs
	c.config.systemKernel = info.SystemKernel
	c.config.systemKernelVersion = info.SystemKernelVersion
	c.config.systemUpTime = info.SystemUpTime
	c.config.systemBootTime = info.SystemBootTime
	c.config.cpuCores = info.CpuCores
	c.config.cpuModelName = info.CpuModelName
	c.config.cpuMhz = info.CpuMhz

	c.config.systemInsideIP = gorequest.GetInsideIp(ctx)
	c.config.systemOutsideIP = systemOutsideIP

	c.config.goVersion = runtime.Version()      // go版本
	c.config.sdkVersion = Version               // sdk版本
	c.config.systemVersion = info.SystemVersion // 系统版本
	c.config.logVersion = golog.Version         // log版本
	c.config.redisSdkVersion = redis.Version()  // redisSdk版本

}

// ConfigGormClientFun GORM配置
func (c *Client) ConfigGormClientFun(ctx context.Context, client *gorm.DB, taskTableName string, taskLogStatus bool, taskLogTableName string) error {
	if client == nil {
		return errors.New("请配置 Gorm")
	}

	// 配置数据库
	c.gormConfig.client = client
	if taskTableName == "" {
		c.gormConfig.taskTableName = "task"
	} else {
		c.gormConfig.taskTableName = taskTableName
	}
	c.gormConfig.taskLogStatus = taskLogStatus
	if c.gormConfig.taskLogStatus {
		if taskLogTableName == "" {
			c.gormConfig.taskLogTableName = "task_log"
		} else {
			c.gormConfig.taskLogTableName = taskLogTableName
		}
	}

	err := c.gormAutoMigrateTask(ctx)
	if err != nil {
		return err
	}
	err = c.gormAutoMigrateTaskLog(ctx)

	return err
}

// ConfigMongoClientFun MONGO配置
func (c *Client) ConfigMongoClientFun(ctx context.Context, client *mongo.Client, databaseName string, taskLogStatus bool, taskLogCollectionName string) error {
	if client == nil {
		return errors.New("请配置 Mongo")
	}

	// 配置数据库
	c.mongoConfig.client = client
	if databaseName == "" {
		return errors.New("请配置 Mongo 库名")
	} else {
		c.mongoConfig.databaseName = databaseName
	}
	c.mongoConfig.taskLogStatus = taskLogStatus
	if c.mongoConfig.taskLogStatus {
		if taskLogCollectionName == "" {
			return errors.New("请配置 Mongo 任务日志集合名")
		} else {
			c.mongoConfig.taskLogCollectionName = taskLogCollectionName
		}
	}

	return nil
}

// ConfigRedisClientFun REDIS配置
// lockKeyPrefix 锁Key前缀 xxx_lock
// lockKeySeparator 锁Key分隔符 :
// cornKeyPrefix 任务Key前缀 xxx_cron
// cornKeyCustom 任务Key自定义 xxx_cron_自定义  xxx_cron_自定义_*
func (c *Client) ConfigRedisClientFun(ctx context.Context, client *redis.Client, lockKeyPrefix string, lockKeySeparator string, cornKeyPrefix string, cornKeyCustom string) error {
	if client == nil {
		return errors.New("请配置 Redis")
	}

	// 配置缓存
	c.redisConfig.client = client

	// 配置缓存前缀
	c.redisConfig.lockKeyPrefix, c.redisConfig.lockKeySeparator, c.redisConfig.cornKeyPrefix, c.redisConfig.cornKeyCustom = lockKeyPrefix, lockKeySeparator, cornKeyPrefix, cornKeyCustom
	if c.redisConfig.lockKeyPrefix == "" || c.redisConfig.lockKeySeparator == "" || c.redisConfig.cornKeyPrefix == "" || c.redisConfig.cornKeyCustom == "" {
		return errors.New("请配置 Redis 前缀")
	}

	return nil
}

// ConfigSLogClientFun 日志配置
func (c *Client) ConfigSLogClientFun(sLogFun golog.SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		c.slog.client = sLog
		c.slog.status = true
	}
}
