package gojobs

import (
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"gorm.io/gorm"
)

// RefreshIp 刷新Ip
func (j *jobsGorm) RefreshIp(tx *gorm.DB) {
	xip := goip.GetOutsideIp()
	if j.outsideIp == "" || j.outsideIp == "0.0.0.0" {
		return
	}
	if j.outsideIp == xip {
		return
	}
	tx.Where("ips = ?", j.outsideIp).Delete(&jobs_gorm_model.TaskIp{}) // 删除
	j.outsideIp = xip
}
