package gojobs

import (
	"context"
	"fmt"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"go.dtapp.net/gotime"
	"go.dtapp.net/gotrace_id"
	"strings"
)

// Filter 过滤
// ctx 上下文
// isMandatoryIp 强制当前ip
// specifyIp 指定Ip
// tasks 过滤前的数据
// newTasks 过滤后的数据
func (c *Client) Filter(ctx context.Context, isMandatoryIp bool, specifyIp string, tasks []jobs_gorm_model.Task, isPrint bool) (newTasks []jobs_gorm_model.Task) {
	c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】是强制性Ip：%v；指定Ip：%v；任务数量：%v", isMandatoryIp, specifyIp, len(tasks)))
	if specifyIp == "" {
		specifyIp = goip.IsIp(c.GetCurrentIp())
	} else {
		specifyIp = goip.IsIp(specifyIp)
	}
	c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】指定Ip重新解析：%v", specifyIp))
	for _, v := range tasks {
		c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】任务指定Ip解析前：%v", v.SpecifyIP))
		v.SpecifyIP = goip.IsIp(v.SpecifyIP)
		c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】任务指定Ip重新解析：%v", v.SpecifyIP))
		// 强制只能是当前的ip
		if isMandatoryIp {
			c.Println(ctx, isPrint, "【Filter入参】进入强制性Ip")
			if v.SpecifyIP == specifyIp {
				c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】进入强制性Ip 添加任务：%v", v.ID))
				newTasks = append(newTasks, v)
				continue
			}
		}
		if v.SpecifyIP == "" {
			c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】任务指定Ip为空 添加任务：%v", v.ID))
			newTasks = append(newTasks, v)
			continue
		} else if v.SpecifyIP == SpecifyIpNull {
			c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】任务指定Ip无限制 添加任务：%v", v.ID))
			newTasks = append(newTasks, v)
			continue
		} else {
			// 判断是否包含该ip
			specifyIpFind := strings.Contains(v.SpecifyIP, ",")
			if specifyIpFind {
				c.Println(ctx, isPrint, "【Filter入参】进入强制性多Ip")
				// 分割字符串
				parts := strings.Split(v.SpecifyIP, ",")
				for _, vv := range parts {
					if vv == specifyIp {
						c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】进入强制性多Ip 添加任务：%v", v.ID))
						newTasks = append(newTasks, v)
						continue
					}
				}
			} else {
				c.Println(ctx, isPrint, "【Filter入参】进入强制性单Ip")
				if v.SpecifyIP == specifyIp {
					newTasks = append(newTasks, v)
					c.Println(ctx, isPrint, fmt.Sprintf("【Filter入参】进入强制性单Ip 添加任务：%v", v.ID))
					continue
				}
			}
		}
	}
	return newTasks
}

// Run 运行
func (c *Client) Run(ctx context.Context, task jobs_gorm_model.Task, taskResultCode int, taskResultDesc string) {

	runId := gotrace_id.GetTraceIdContext(ctx)
	if runId == "" {
		if c.slog.status {
			c.slog.client.WithTraceId(ctx).Error("上下文没有跟踪编号")
		}
		return
	}

	c.TaskLogRecord(ctx, task, runId, taskResultCode, taskResultDesc)

	switch taskResultCode {
	case 0:
		err := c.EditTask(ctx, c.gormClient, task.ID).
			Select("run_id", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				RunID:       runId,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().AfterSeconds(task.Frequency).Time,
			}).Error
		if err != nil {
			if c.slog.status {
				c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("保存失败：%s", err))
			}
		}
		return
	case CodeSuccess:
		// 执行成功
		err := c.EditTask(ctx, c.gormClient, task.ID).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行成功",
				Number:      task.Number + 1,
				RunID:       runId,
				UpdatedIP:   c.config.systemOutsideIP,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().AfterSeconds(task.Frequency).Time,
			}).Error
		if err != nil {
			if c.slog.status {
				c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("保存失败：%s", err))
			}
		}
	case CodeEnd:
		// 执行成功、提前结束
		err := c.EditTask(ctx, c.gormClient, task.ID).
			Select("status", "status_desc", "number", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				Status:      TASK_SUCCESS,
				StatusDesc:  "结束执行",
				Number:      task.Number + 1,
				UpdatedIP:   c.config.systemOutsideIP,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().Time,
			}).Error
		if err != nil {
			if c.slog.status {
				c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("保存失败：%s", err))
			}
		}
	case CodeError:
		// 执行失败
		err := c.EditTask(ctx, c.gormClient, task.ID).
			Select("status_desc", "number", "run_id", "updated_ip", "result", "next_run_time").
			Updates(jobs_gorm_model.Task{
				StatusDesc:  "执行失败",
				Number:      task.Number + 1,
				RunID:       runId,
				UpdatedIP:   c.config.systemOutsideIP,
				Result:      taskResultDesc,
				NextRunTime: gotime.Current().AfterSeconds(task.Frequency).Time,
			}).Error
		if err != nil {
			if c.slog.status {
				c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("保存失败：%s", err))
			}
		}
	}

	if task.MaxNumber != 0 {
		if task.Number+1 >= task.MaxNumber {
			// 关闭执行
			err := c.EditTask(ctx, c.gormClient, task.ID).
				Select("status").
				Updates(jobs_gorm_model.Task{
					Status: TASK_TIMEOUT,
				}).Error
			if err != nil {
				if c.slog.status {
					c.slog.client.WithTraceId(ctx).Error(fmt.Sprintf("保存失败：%s", err))
				}
			}
		}
	}
	return
}
