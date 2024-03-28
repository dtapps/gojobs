package gorequest

import (
	"github.com/shirou/gopsutil/host"
	"runtime"
)

type systemResult struct {
	SystemOs     string // 系统类型
	SystemKernel string // 系统内核
}

// 获取系统信息
func getSystem() (result systemResult) {
	hInfo, _ := host.Info()
	result.SystemOs = hInfo.OS
	result.SystemKernel = hInfo.KernelArch
	return result
}

// 设置配置信息
func (app *App) setConfig() {
	info := getSystem()
	app.config.systemOs = info.SystemOs
	app.config.systemKernel = info.SystemKernel
	app.config.goVersion = runtime.Version()
	app.config.sdkVersion = Version
}
