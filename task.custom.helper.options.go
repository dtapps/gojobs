package gojobs

type TaskCustomHelperOption interface {
	apply(cfg *taskCustomHelperConfig)
}

type taskCustomHelperOption func(cfg *taskCustomHelperConfig)

func (fn taskCustomHelperOption) apply(cfg *taskCustomHelperConfig) {
	fn(cfg)
}

type taskCustomHelperConfig struct {
	logIsDebug            bool   // [日志]日志是否启动
	traceIsFilter         bool   // [过滤]链路追踪是否过滤
	traceIsFilterKeyName  string // [过滤]Key名称
	traceIsFilterKeyValue string // [过滤]Key值
}

// defaultTaskCustomHelperConfig 默认配置
func defaultTaskCustomHelperConfig() *taskCustomHelperConfig {
	return &taskCustomHelperConfig{}
}

// newTaskCustomHelperConfig 初始配置
func newTaskCustomHelperConfig(opts []TaskCustomHelperOption) *taskCustomHelperConfig {
	cfg := defaultTaskCustomHelperConfig()
	for _, opt := range opts {
		opt.apply(cfg)
	}
	return cfg
}

// TaskCustomHelperWithDebug 设置日志是否打印
func TaskCustomHelperWithDebug(is bool) TaskCustomHelperOption {
	return taskCustomHelperOption(func(cfg *taskCustomHelperConfig) {
		cfg.logIsDebug = is
	})
}

// TaskCustomHelperWithFilter 设置链路追踪是否过滤
func TaskCustomHelperWithFilter(is bool, keyName string, keyValue string) TaskCustomHelperOption {
	return taskCustomHelperOption(func(cfg *taskCustomHelperConfig) {
		cfg.traceIsFilter = is
		cfg.traceIsFilterKeyName = keyName
		cfg.traceIsFilterKeyValue = keyValue
	})
}
