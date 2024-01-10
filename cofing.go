package gojobs

import (
	"context"
	"github.com/redis/go-redis/v9"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"go.dtapp.net/golog"
	"go.dtapp.net/gorequest"
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
func (c *Client) setConfig(ctx context.Context, systemOutsideIp string) {

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
	c.config.systemOutsideIP = systemOutsideIp

	c.config.goVersion = runtime.Version()      // go版本
	c.config.sdkVersion = Version               // sdk版本
	c.config.systemVersion = info.SystemVersion // 系统版本
	c.config.logVersion = golog.Version         // log版本
	c.config.redisSdkVersion = redis.Version()  // redis版本

}

// ConfigSLogClientFun 日志配置
func (c *Client) ConfigSLogClientFun(sLogFun golog.SLogFun) {
	sLog := sLogFun()
	if sLog != nil {
		c.slog.client = sLog
		c.slog.status = true
	}
}
