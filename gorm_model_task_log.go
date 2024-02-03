package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gotime"
	"log"
	"time"
)

// GormModelTaskLog 任务日志
type GormModelTaskLog struct {
	LogID           uint      `gorm:"primaryKey;comment:【日志】编号" json:"log_id"`                            // 【日志】编号
	LogTime         time.Time `gorm:"autoCreateTime;index;comment:【日志】时间" json:"log_time"`                // 【日志】时间
	TaskID          uint      `gorm:"index;comment:【任务】编号" json:"task_id"`                                // 【任务】编号
	TaskRunID       string    `gorm:"comment:【任务】执行编号" json:"task_run_id"`                                //【任务】执行编号
	TaskResultCode  int       `gorm:"index;comment:【任务】执行状态码" json:"task_result_code"`                    //【任务】执行状态码
	TaskResultDesc  string    `gorm:"comment:【任务】执行结果" json:"task_result_desc"`                           //【任务】执行结果
	SystemHostName  string    `gorm:"comment:【系统】主机名" json:"system_host_name,omitempty"`                  //【系统】主机名
	SystemInsideIP  string    `gorm:"default:0.0.0.0;comment:【系统】内网IP" json:"system_inside_ip,omitempty"` //【系统】内网IP
	SystemOutsideIP string    `gorm:"default:0.0.0.0;comment:【系统】外网IP" json:"system_outside_ip"`          //【系统】外网IP
	SystemOs        string    `gorm:"comment:【系统】类型" json:"system_os,omitempty"`                          //【系统】类型
	SystemArch      string    `gorm:"comment:【系统】架构" json:"system_arch,omitempty"`                        //【系统】架构
	SystemUpTime    uint64    `gorm:"comment:【系统】运行时间" json:"system_up_time,omitempty"`                   //【系统】运行时间
	SystemBootTime  uint64    `gorm:"comment:【系统】开机时间" json:"system_boot_time,omitempty"`                 //【系统】开机时间
	GoVersion       string    `gorm:"comment:【程序】Go版本" json:"go_version,omitempty"`                       //【程序】Go版本
	SdkVersion      string    `gorm:"comment:【程序】Sdk版本" json:"sdk_version,omitempty"`                     //【程序】Sdk版本
	SystemVersion   string    `gorm:"comment:【程序】System版本" json:"system_version,omitempty"`               //【程序】System版本
	CpuCores        int       `gorm:"comment:【CPU】核数" json:"cpu_cores,omitempty"`                         //【CPU】核数
	CpuModelName    string    `gorm:"comment:【CPU】型号名称" json:"cpu_model_name,omitempty"`                  //【CPU】型号名称
	CpuMhz          float64   `gorm:"comment:【CPU】兆赫" json:"cpu_mhz,omitempty"`                           //【CPU】兆赫
}

// 创建模型
func (c *Client) gormAutoMigrateTaskLog(ctx context.Context) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		AutoMigrate(&GormModelTaskLog{})
	return err
}

// GormTaskLogDelete 删除
func (c *Client) GormTaskLogDelete(ctx context.Context, hour int64) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&GormModelTaskLog{}).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogInDelete 删除任务运行
func (c *Client) GormTaskLogInDelete(ctx context.Context, hour int64) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Where("task_result_status = ?", TASK_IN).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&GormModelTaskLog{}).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogSuccessDelete 删除任务完成
func (c *Client) GormTaskLogSuccessDelete(ctx context.Context, hour int64) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Where("task_result_status = ?", TASK_SUCCESS).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&GormModelTaskLog{}).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogErrorDelete 删除任务异常
func (c *Client) GormTaskLogErrorDelete(ctx context.Context, hour int64) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Where("task_result_status = ?", TASK_ERROR).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&GormModelTaskLog{}).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogTimeoutDelete 删除任务超时
func (c *Client) GormTaskLogTimeoutDelete(ctx context.Context, hour int64) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Where("task_result_status = ?", TASK_TIMEOUT).Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&GormModelTaskLog{}).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogWaitDelete 删除任务等待
func (c *Client) GormTaskLogWaitDelete(ctx context.Context, hour int64) error {
	if c.gormConfig.taskLogStatus == false {
		return nil
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Where("task_result_status = ?", TASK_WAIT).
		Where("log_time < ?", gotime.Current().BeforeHour(hour).Format()).
		Delete(&GormModelTaskLog{}).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("删除失败：%s", err))
		}
	}
	return err
}

// GormTaskLogRecord 记录
func (c *Client) GormTaskLogRecord(ctx context.Context, task GormModelTask, runId string, taskResultCode int, taskResultDesc string) {
	taskLog := GormModelTaskLog{
		TaskID:          task.ID,                  //【任务】编号
		TaskRunID:       runId,                    //【任务】执行编号
		TaskResultCode:  taskResultCode,           //【任务】执行状态码
		TaskResultDesc:  taskResultDesc,           //【任务】执行结果
		SystemHostName:  c.config.systemHostname,  //【系统】主机名
		SystemInsideIP:  c.config.systemInsideIP,  //【系统】内网IP
		SystemOutsideIP: c.config.systemOutsideIP, //【系统】外网IP
		SystemOs:        c.config.systemOs,        //【系统】类型
		SystemArch:      c.config.systemKernel,    //【系统】架构
		SystemUpTime:    c.config.systemUpTime,    //【系统】运行时间
		SystemBootTime:  c.config.systemBootTime,  //【系统】开机时间
		GoVersion:       c.config.goVersion,       //【程序】Go版本
		SdkVersion:      c.config.sdkVersion,      //【程序】Sdk版本
		SystemVersion:   c.config.sdkVersion,      //【程序】System版本
		CpuCores:        c.config.cpuCores,        //【程序】核数
		CpuModelName:    c.config.cpuModelName,    //【程序】型号名称
		CpuMhz:          c.config.cpuMhz,          //【程序】兆赫
	}
	err := c.gormConfig.client.WithContext(ctx).Table(c.gormConfig.taskLogTableName).
		Create(&taskLog).Error
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("记录失败：%s", err))
		}
		if c.slog.status {
			log.Println(fmt.Sprintf("记录数据：%+v", taskLog))
		}
	}
}
