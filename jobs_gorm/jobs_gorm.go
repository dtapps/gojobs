package jobs_gorm

import "gorm.io/gorm"
import "go.dtapp.net/goredis"

type JobsGorm struct {
	RunVersion  string      // 运行版本
	Os          string      // 系统类型
	Arch        string      // 系统架构
	MaxProCs    int         // CPU核数
	Version     string      // GO版本
	MacAddrS    string      // Mac地址
	InsideIp    string      // 内网ip
	OutsideIp   string      // 外网ip
	MainService int         // 主要服务
	Db          *gorm.DB    // 数据库
	Redis       goredis.App // 缓存数据库服务
}
