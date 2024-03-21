package gojobs

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.dtapp.net/gotime"
	"log"
)

type taskList struct {
	id   cron.EntryID
	name string
}

// Cron 定时任务管理器
type Cron struct {
	inner *cron.Cron
	list  []taskList
}

// NewCron 创建一个定时任务管理器
func NewCron() *Cron {
	return &Cron{
		inner: cron.New(),
		list:  make([]taskList, 0),
	}
}

func NewCronWithSeconds() *Cron {
	return &Cron{
		inner: cron.New(cron.WithSeconds()),
		list:  make([]taskList, 0),
	}
}

// Start 启动任务
func (c *Cron) Start() {
	c.inner.Start()
}

// Stop 关闭任务
func (c *Cron) Stop() context.Context {
	return c.inner.Stop()
}

// AddFunc 添加任务
func (c *Cron) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	id, err := c.inner.AddFunc(spec, cmd)
	c.list = append(c.list, taskList{
		id: id,
	})
	return id, err
}

// AddJob 添加任务
func (c *Cron) AddJob(spec string, cmd cron.Job) (cron.EntryID, error) {
	id, err := c.inner.AddJob(spec, cmd)
	c.list = append(c.list, taskList{
		id: id,
	})
	return id, err
}

// Entry 查询任务
func (c *Cron) Entry(id cron.EntryID) cron.Entry {
	return c.inner.Entry(id)
}

// Remove 删除任务
func (c *Cron) Remove(id cron.EntryID) {
	c.inner.Remove(id)
}

// List 任务列表
func (c *Cron) List() []cron.EntryID {
	ids := make([]cron.EntryID, 0)
	for _, v := range c.list {
		ids = append(ids, v.id)
	}
	return ids
}

// ListShow 任务列表
func (c *Cron) ListShow() {
	for _, v := range c.list {
		taskInfo := c.inner.Entry(v.id)
		log.Println(fmt.Sprintf("[ID=%v][Schedule=%v][Prev=%v][Next=%v]",
			taskInfo.ID,
			taskInfo.Schedule,
			taskInfo.Prev.Format(gotime.DateTimeZhFormat),
			taskInfo.Next.Format(gotime.DateTimeZhFormat),
		))
	}
}

// AddTask 添加任务
func (c *Cron) AddTask(name string, spec string, cmd func()) (cron.EntryID, error) {
	id, err := c.inner.AddFunc(spec, cmd)
	c.list = append(c.list, taskList{
		id:   id,
		name: name,
	})
	return id, err
}

// QueryTask 查询任务
func (c *Cron) QueryTask(id cron.EntryID) cron.Entry {
	return c.inner.Entry(id)
}

// RemoveTask 删除任务
func (c *Cron) RemoveTask(id cron.EntryID) {
	c.inner.Remove(id)
}

// ListTask 任务列表
func (c *Cron) ListTask() {
	for _, v := range c.list {
		taskInfo := c.inner.Entry(v.id)
		log.Println(fmt.Sprintf("%s [ID=%v][Schedule=%v][Prev=%v][Next=%v]",
			v.name,
			taskInfo.ID,
			taskInfo.Schedule,
			taskInfo.Prev.Format(gotime.DateTimeZhFormat),
			taskInfo.Next.Format(gotime.DateTimeZhFormat),
		))
	}
}
