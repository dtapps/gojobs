package gojobs

import (
	"context"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_gorm_model"
	"gorm.io/gorm"
)

// RefreshIp 刷新Ip
func (c *Client) RefreshIp(ctx context.Context, tx *gorm.DB) {
	xip := goip.GetOutsideIp(ctx)
	if c.config.outsideIp == "" || c.config.outsideIp == "0.0.0.0" {
		return
	}
	if c.config.outsideIp == xip {
		return
	}
	tx.Where("ips = ?", c.config.outsideIp).Delete(&jobs_gorm_model.TaskIp{}) // 删除
	c.config.outsideIp = xip
}
