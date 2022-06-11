package jobs_gorm

import (
	"go.dtapp.net/goip"
	"gorm.io/gorm"
)

// RefreshIp 刷新Ip
func (jobsGorm *JobsGorm) RefreshIp(tx *gorm.DB) {
	xip := goip.GetOutsideIp()
	if jobsGorm.outsideIp == "" || jobsGorm.outsideIp == "0.0.0.0" {
		return
	}
	if jobsGorm.outsideIp == xip {
		return
	}
	tx.Where("ips = ?", jobsGorm.outsideIp).Delete(&TaskIp{}) // 删除
	jobsGorm.outsideIp = xip
}
