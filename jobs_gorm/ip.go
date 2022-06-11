package jobs_gorm

import (
	"go.dtapp.net/goip"
	"gorm.io/gorm"
)

// RefreshIp 刷新Ip
func (jobsGorm *JobsGorm) RefreshIp(tx *gorm.DB) {
	xip := goip.GetOutsideIp()
	if jobsGorm.OutsideIp == "" || jobsGorm.OutsideIp == "0.0.0.0" {
		return
	}
	if jobsGorm.OutsideIp == xip {
		return
	}
	tx.Where("ips = ?", jobsGorm.OutsideIp).Delete(&TaskIp{}) // 删除
	jobsGorm.OutsideIp = xip
}
