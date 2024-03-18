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

func (c *Cron) AddFunc(spec string, cmd func()) (cron.EntryID, error) {
	id, err := c.inner.AddFunc(spec, cmd)
	c.ids = append(c.ids, id)
	return id, err
}

func (c *Cron) List() []cron.EntryID {
	return c.ids
}
func (c *Cron) ListShow() {
	for _, v := range c.ids {
		taskInfo := c.inner.Entry(v)
		log.Println(fmt.Sprintf("[ID=%v][Schedule=%v]{Prev=%v}{Next=%v}",
			taskInfo.ID,
			taskInfo.Schedule,
			taskInfo.Prev.Format(gotime.DateTimeZhFormat),
			taskInfo.Next.Format(gotime.DateTimeZhFormat),
		))
	}
}
