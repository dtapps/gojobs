package gojobs

import (
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"go.dtapp.net/gotime"
	"log"
)

// Cron 定时任务管理器
type Cron struct {
	inner *cron.Cron
	ids   []cron.EntryID
}

// NewCron 创建一个定时任务管理器
func NewCron() *Cron {
	return &Cron{
		inner: cron.New(),
		ids:   make([]cron.EntryID, 0),
	}
}

func NewCronWithSeconds() *Cron {
	return &Cron{
		inner: cron.New(cron.WithSeconds()),
		ids:   make([]cron.EntryID, 0),
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
	c.ids = append(c.ids, id)
	return id, err
}

// AddJob 添加任务
func (c *Cron) AddJob(spec string, cmd cron.Job) (cron.EntryID, error) {
	id, err := c.inner.AddJob(spec, cmd)
	c.ids = append(c.ids, id)
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
	return c.ids
}

// ListShow 任务列表
func (c *Cron) ListShow() {
	for _, v := range c.ids {
		taskInfo := c.inner.Entry(v)
		log.Println(fmt.Sprintf("[ID=%v][Schedule=%v][Prev=%v][Next=%v]",
			taskInfo.ID,
			taskInfo.Schedule,
			taskInfo.Prev.Format(gotime.DateTimeZhFormat),
			taskInfo.Next.Format(gotime.DateTimeZhFormat),
		))
	}
}
