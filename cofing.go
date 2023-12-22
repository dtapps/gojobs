package gojobs

import (
	"go.dtapp.net/golog"
)

// ConfigSLogClientFun 日志配置
func (c *Client) ConfigSLogClientFun(sLogFun golog.SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		c.slog.client = sLog
		c.slog.status = true
	}
}

// ConfigRunSLogClientFun 运行日志配置
func (c *Client) ConfigRunSLogClientFun(sLogFun golog.SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		c.runSlog.client = sLog
		c.runSlog.status = true
	}
}
