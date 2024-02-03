package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/gotime"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

// MongoModelTaskLog 任务日志
type MongoModelTaskLog struct {
	LogID           primitive.ObjectID `json:"log_id,omitempty" bson:"_id,omitempty"`                        //【记录】编号
	LogTime         primitive.DateTime `json:"log_time,omitempty" bson:"log_time"`                           //【记录】时间
	TaskID          uint               `json:"task_id" bson:"task_id,omitempty"`                             //【任务】编号
	TaskRunID       string             `json:"task_run_id" bson:"task_run_id,omitempty"`                     //【任务】执行编号
	TaskResultCode  int                `json:"task_result_code" bson:"task_result_code,omitempty"`           //【任务】执行状态码
	TaskResultDesc  string             `json:"task_result_desc" bson:"task_result_desc,omitempty"`           //【任务】执行结果
	SystemHostName  string             `json:"system_host_name,omitempty" bson:"system_host_name,omitempty"` //【系统】主机名
	SystemInsideIP  string             `json:"system_inside_ip,omitempty" bson:"system_inside_ip,omitempty"` //【系统】内网IP
	SystemOutsideIP string             `json:"system_outside_ip" bson:"system_outside_ip,omitempty"`         //【系统】外网IP
	SystemOs        string             `json:"system_os,omitempty" bson:"system_os,omitempty"`               //【系统】类型
	SystemArch      string             `json:"system_arch,omitempty" bson:"system_arch,omitempty"`           //【系统】架构
	SystemUpTime    uint64             `json:"system_up_time,omitempty" bson:"system_up_time,omitempty"`     //【系统】运行时间
	SystemBootTime  uint64             `json:"system_boot_time,omitempty" bson:"system_boot_time,omitempty"` //【系统】开机时间
	GoVersion       string             `json:"go_version,omitempty" bson:"go_version,omitempty"`             //【程序】Go版本
	SdkVersion      string             `json:"sdk_version,omitempty" bson:"sdk_version,omitempty"`           //【程序】Sdk版本
	SystemVersion   string             `json:"system_version,omitempty" bson:"system_version,omitempty"`     //【程序】System版本
	CpuCores        int                `json:"cpu_cores,omitempty" bson:"cpu_cores,omitempty"`               //【CPU】核数
	CpuModelName    string             `json:"cpu_model_name,omitempty" bson:"cpu_model_name,omitempty"`     //【CPU】型号名称
	CpuMhz          float64            `json:"cpu_mhz,omitempty" bson:"cpu_mhz,omitempty"`                   //【CPU】兆赫
}

// 创建时间序列集合
func (c *Client) mongoCreateCollectionTaskLog(ctx context.Context) error {
	if c.mongoConfig.taskLogStatus == false {
		return nil
	}
	err := c.mongoConfig.client.Database(c.mongoConfig.databaseName).
		CreateCollection(ctx,
			c.mongoConfig.taskLogCollectionName,
			options.CreateCollection().SetTimeSeriesOptions(options.TimeSeries().SetTimeField("log_time")))
	return err
}

// 创建索引
func (c *Client) mongoCreateIndexesTaskLog(ctx context.Context) error {
	if c.mongoConfig.taskLogStatus == false {
		return nil
	}
	_, err := c.mongoConfig.client.Database(c.mongoConfig.databaseName).
		Collection(c.mongoConfig.taskLogCollectionName).
		Indexes().
		CreateMany(ctx, []mongo.IndexModel{
			{
				Keys: bson.D{{
					Key:   "log_time",
					Value: -1,
				}},
			}})
	return err
}

// MongoTaskLogRecord 记录
func (c *Client) MongoTaskLogRecord(ctx context.Context, task GormModelTask, runId string, taskResultCode int, taskResultDesc string) {
	taskLog := MongoModelTaskLog{
		LogID:           primitive.NewObjectID(),                              //【记录】编号
		LogTime:         primitive.NewDateTimeFromTime(gotime.Current().Time), //【记录】时间
		TaskID:          task.ID,                                              //【任务】编号
		TaskRunID:       runId,                                                //【任务】执行编号
		TaskResultCode:  taskResultCode,                                       //【任务】执行状态码
		TaskResultDesc:  taskResultDesc,                                       //【任务】执行结果
		SystemHostName:  c.config.systemHostname,                              //【系统】主机名
		SystemInsideIP:  c.config.systemInsideIP,                              //【系统】内网IP
		SystemOutsideIP: c.config.systemOutsideIP,                             //【系统】外网IP
		SystemOs:        c.config.systemOs,                                    //【系统】类型
		SystemArch:      c.config.systemKernel,                                //【系统】架构
		SystemUpTime:    c.config.systemUpTime,                                //【系统】运行时间
		SystemBootTime:  c.config.systemBootTime,                              //【系统】开机时间
		GoVersion:       c.config.goVersion,                                   //【程序】Go版本
		SdkVersion:      c.config.sdkVersion,                                  //【程序】Sdk版本
		SystemVersion:   c.config.sdkVersion,                                  //【程序】System版本
		CpuCores:        c.config.cpuCores,                                    //【程序】核数
		CpuModelName:    c.config.cpuModelName,                                //【程序】型号名称
		CpuMhz:          c.config.cpuMhz,                                      //【程序】兆赫
	}
	_, err := c.mongoConfig.client.Database(c.mongoConfig.databaseName).
		Collection(c.mongoConfig.taskLogCollectionName).
		InsertOne(ctx, taskLog)
	if err != nil {
		if c.slog.status {
			log.Println(fmt.Sprintf("记录失败：%s", err))
		}
		if c.slog.status {
			log.Println(fmt.Sprintf("记录数据：%+v", taskLog))
		}
	}
}
