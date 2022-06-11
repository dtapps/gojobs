package jobs_gorm

import (
	"errors"
	"fmt"
	"go.dtapp.net/goarray"
	"go.dtapp.net/goip"
	"go.dtapp.net/gojobs/jobs_common"
	"go.dtapp.net/goredis"
	"go.dtapp.net/gouuid"
	"gorm.io/gorm"
	"runtime"
)

// JobsGorm 任务
type JobsGorm struct {
	runVersion  string          // 运行版本
	os          string          // 系统类型
	arch        string          // 系统架构
	maxProCs    int             // CPU核数
	version     string          // GO版本
	macAddrS    string          // Mac地址
	insideIp    string          // 内网ip
	outsideIp   string          // 外网ip
	mainService int             // 主要服务
	Db          *gorm.DB        // 数据库
	Redis       *goredis.Client // 缓存数据库服务
}

// NewGorm 任务
func NewGorm(jobsGorm JobsGorm, mainService int, runVersion string) *JobsGorm {
	jobsGorm.runVersion = runVersion
	jobsGorm.os = runtime.GOOS
	jobsGorm.arch = runtime.GOARCH
	jobsGorm.maxProCs = runtime.GOMAXPROCS(0)
	jobsGorm.version = runtime.Version()
	jobsGorm.macAddrS = goarray.TurnString(goip.GetMacAddr())
	jobsGorm.insideIp = goip.GetInsideIp()
	jobsGorm.outsideIp = goip.GetOutsideIp()
	jobsGorm.mainService = mainService
	return &jobsGorm
}

// ConfigCreateInCustomId 创建正在运行任务
type ConfigCreateInCustomId struct {
	Tx             *gorm.DB // 驱动
	Params         string   // 参数
	Frequency      int64    // 频率(秒单位)
	CustomId       string   // 自定义编号
	CustomSequence int64    // 自定义顺序
	Type           string   // 类型
	SpecifyIp      string   // 指定外网IP
}

// CreateInCustomId 创建正在运行任务
func (jobsGorm *JobsGorm) CreateInCustomId(config *ConfigCreateInCustomId) error {
	createStatus := config.Tx.Create(&Task{
		Status:         jobs_common.TASK_IN,
		Params:         config.Params,
		StatusDesc:     "首次添加任务",
		Frequency:      config.Frequency,
		RunId:          gouuid.GetUuId(),
		CustomId:       config.CustomId,
		CustomSequence: config.CustomSequence,
		Type:           config.Type,
		CreatedIp:      jobsGorm.outsideIp,
		SpecifyIp:      config.SpecifyIp,
		UpdatedIp:      jobsGorm.outsideIp,
	})
	if createStatus.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("创建任务失败：%s", createStatus.Error))
	}
	return nil
}

// ConfigCreateInCustomIdOnly 创建正在运行唯一任务
type ConfigCreateInCustomIdOnly struct {
	Tx             *gorm.DB // 驱动
	Params         string   // 参数
	Frequency      int64    // 频率(秒单位)
	CustomId       string   // 自定义编号
	CustomSequence int64    // 自定义顺序
	Type           string   // 类型
	SpecifyIp      string   // 指定外网IP
}

// CreateInCustomIdOnly 创建正在运行唯一任务
func (jobsGorm *JobsGorm) CreateInCustomIdOnly(config *ConfigCreateInCustomIdOnly) error {
	query := jobsGorm.TaskTypeTakeIn(config.Tx, config.CustomId, config.Type)
	if query.Id != 0 {
		return errors.New("任务已存在")
	}
	createStatus := config.Tx.Create(&Task{
		Status:         jobs_common.TASK_IN,
		Params:         config.Params,
		StatusDesc:     "首次添加任务",
		Frequency:      config.Frequency,
		RunId:          gouuid.GetUuId(),
		CustomId:       config.CustomId,
		CustomSequence: config.CustomSequence,
		Type:           config.Type,
		CreatedIp:      jobsGorm.outsideIp,
		SpecifyIp:      config.SpecifyIp,
		UpdatedIp:      jobsGorm.outsideIp,
	})
	if createStatus.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("创建任务失败：%s", createStatus.Error))
	}
	return nil
}

// ConfigCreateInCustomIdMaxNumber 创建正在运行任务并限制数量
type ConfigCreateInCustomIdMaxNumber struct {
	Tx             *gorm.DB // 驱动
	Params         string   // 参数
	Frequency      int64    // 频率(秒单位)
	MaxNumber      int64    // 最大次数
	CustomId       string   // 自定义编号
	CustomSequence int64    // 自定义顺序
	Type           string   // 类型
	SpecifyIp      string   // 指定外网IP
}

// CreateInCustomIdMaxNumber 创建正在运行任务并限制数量
func (jobsGorm *JobsGorm) CreateInCustomIdMaxNumber(config *ConfigCreateInCustomIdMaxNumber) error {
	createStatus := config.Tx.Create(&Task{
		Status:         jobs_common.TASK_IN,
		Params:         config.Params,
		StatusDesc:     "首次添加任务",
		Frequency:      config.Frequency,
		MaxNumber:      config.MaxNumber,
		RunId:          gouuid.GetUuId(),
		CustomId:       config.CustomId,
		CustomSequence: config.CustomSequence,
		Type:           config.Type,
		CreatedIp:      jobsGorm.outsideIp,
		SpecifyIp:      config.SpecifyIp,
		UpdatedIp:      jobsGorm.outsideIp,
	})
	if createStatus.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("创建任务失败：%s", createStatus.Error))
	}
	return nil
}

// ConfigCreateInCustomIdMaxNumberOnly 创建正在运行唯一任务并限制数量
type ConfigCreateInCustomIdMaxNumberOnly struct {
	Tx             *gorm.DB // 驱动
	Params         string   // 参数
	Frequency      int64    // 频率(秒单位)
	MaxNumber      int64    // 最大次数
	CustomId       string   // 自定义编号
	CustomSequence int64    // 自定义顺序
	Type           string   // 类型
	SpecifyIp      string   // 指定外网IP
}

// CreateInCustomIdMaxNumberOnly 创建正在运行唯一任务并限制数量
func (jobsGorm *JobsGorm) CreateInCustomIdMaxNumberOnly(config *ConfigCreateInCustomIdMaxNumberOnly) error {
	query := jobsGorm.TaskTypeTakeIn(config.Tx, config.CustomId, config.Type)
	if query.Id != 0 {
		return errors.New("任务已存在")
	}
	createStatus := config.Tx.Create(&Task{
		Status:         jobs_common.TASK_IN,
		Params:         config.Params,
		StatusDesc:     "首次添加任务",
		Frequency:      config.Frequency,
		MaxNumber:      config.MaxNumber,
		RunId:          gouuid.GetUuId(),
		CustomId:       config.CustomId,
		CustomSequence: config.CustomSequence,
		Type:           config.Type,
		CreatedIp:      jobsGorm.outsideIp,
		SpecifyIp:      config.SpecifyIp,
		UpdatedIp:      jobsGorm.outsideIp,
	})
	if createStatus.RowsAffected == 0 {
		return errors.New(fmt.Sprintf("创建任务失败：%s", createStatus.Error))
	}
	return nil
}
