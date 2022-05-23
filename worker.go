package goip

import (
	"go.dtapp.net/gojobs/pb"
	"google.golang.org/grpc"
)

// WorkerConfig 工作配置
type WorkerConfig struct {
	Address string // 服务端口 127.0.0.1:8888
}

// Worker 工作
type Worker struct {
	WorkerConfig                  // 配置
	Pub          pb.PubSubClient  // 订阅
	Conn         *grpc.ClientConn // 链接信息
}

// NewWorker 创建工作
func NewWorker(config *WorkerConfig) *Worker {

	if config.Address == "" {
		panic("请填写服务端口")
	}

	w := &Worker{}

	w.Address = config.Address

	var err error

	// 建立连接 获取client
	w.Conn, err = grpc.Dial(w.Address, grpc.WithInsecure())
	if err != nil {
		panic("连接失败: " + err.Error())
	}

	// 新建一个客户端
	w.Pub = pb.NewPubSubClient(w.Conn)

	return w
}
