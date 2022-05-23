package goip

import (
	"go.dtapp.net/gojobs/pb"
	"google.golang.org/grpc"
)

// CronConfig 定时任务配置
type CronConfig struct {
	Address string // 服务端口 127.0.0.1:8888
}

// Cron 定时任务
type Cron struct {
	CronConfig                  // 配置
	Pub        pb.PubSubClient  // 订阅
	Conn       *grpc.ClientConn // 链接信息
}

// NewCron 创建定时任务
func NewCron(config *CronConfig) *Cron {

	if config.Address == "" {
		panic("请填写服务端口")
	}

	c := &Cron{}

	c.Address = config.Address

	var err error

	// 建立连接 获取client
	c.Conn, err = grpc.Dial(c.Address, grpc.WithInsecure())
	if err != nil {
		panic("连接失败: " + err.Error())
	}

	// 新建一个客户端
	c.Pub = pb.NewPubSubClient(c.Conn)

	return c
}
