package gojobs

import (
	"context"
	"github.com/go-redis/redis/v9"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"go.dtapp.net/goip"
	"go.dtapp.net/golog"
	"go.mongodb.org/mongo-driver/version"
	"log"
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

func getSystem() (result systemResult) {

	hInfo, err := host.Info()
	if err != nil {
		log.Printf("getSystem.host.Info：%s\n", err)
	}

	result.SystemHostname = hInfo.Hostname
	result.SystemOs = hInfo.OS
	result.SystemVersion = hInfo.PlatformVersion
	result.SystemKernel = hInfo.KernelArch
	result.SystemKernelVersion = hInfo.KernelVersion
	result.SystemUpTime = hInfo.Uptime
	if hInfo.BootTime != 0 {
		result.SystemBootTime = hInfo.BootTime
	}

	hCpu, err := cpu.Times(true)
	if err != nil {
		log.Printf("getSystem.cpu.Times：%s\n", err)
	}

	result.CpuCores = len(hCpu)

	cInfo, err := cpu.Info()
	if err != nil {
		log.Printf("getSystem.cpu.Info：%s\n", err)
	}
	if len(cInfo) > 0 {
		result.CpuModelName = cInfo[0].ModelName
		result.CpuMhz = cInfo[0].Mhz
	}

	return result
}

func (c *Client) setConfig(ctx context.Context) {

	info := getSystem()

	c.config.systemHostname = info.SystemHostname
	c.config.systemOs = info.SystemOs
	c.config.systemVersion = info.SystemVersion
	c.config.systemKernel = info.SystemKernel
	c.config.systemKernelVersion = info.SystemKernelVersion
	c.config.systemBootTime = info.SystemBootTime
	c.config.cpuCores = info.CpuCores
	c.config.cpuModelName = info.CpuModelName
	c.config.cpuMhz = info.CpuMhz

	c.config.systemInsideIp = goip.GetInsideIp(ctx)

	c.config.sdkVersion = Version
	c.config.goVersion = runtime.Version()

	c.config.mongoSdkVersion = version.Driver

	c.config.redisSdkVersion = redis.Version()

	c.config.logVersion = golog.Version
}
