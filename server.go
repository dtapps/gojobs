package gojobs

import (
	"errors"
	"go.dtapp.net/gojobs/pb"
	"go.dtapp.net/gojobs/pubsub"
	"google.golang.org/grpc"
	"log"
	"net"
	"time"
)

// ServerConfig 服务配置
type ServerConfig struct {
	PublishTimeout time.Duration // 控制发布时最大阻塞时间
	PubBuffer      int           // 缓冲区大小，控制每个订阅者的chan缓冲区大小
	Address        string        // 服务端口 0.0.0.0:8888
}

// Server 服务
type Server struct {
	ServerConfig                   // 配置
	Pub          *pubsub.Publisher // 订阅
	Conn         *grpc.Server      // 链接信息
}

// NewServer 创建服务和注册
func NewServer(config *ServerConfig) *Server {

	if config.Address == "" {
		panic("[服务中转]请填写服务端口")
	}

	s := &Server{}

	s.PublishTimeout = config.PublishTimeout
	s.PubBuffer = config.PubBuffer
	s.Address = config.Address

	s.Pub = pubsub.NewPublisher(config.PublishTimeout, config.PubBuffer)

	// 创建gRPC服务器
	s.Conn = grpc.NewServer()

	// 注册
	pb.RegisterPubSubServer(s.Conn, pb.NewPubSubServerService())

	return s
}

// StartUp 启动
func (s *Server) StartUp() error {

	// 监听本地端口
	lis, err := net.Listen("tcp", s.Address)
	if err != nil {
		return errors.New("[服务中转]{创建监听失败} " + err.Error())
	}
	log.Println("[服务中转]{监听] ", lis.Addr())

	// 启动grpc
	err = s.Conn.Serve(lis)
	if err != nil {
		return errors.New("[服务中转]{创建服务失败} " + err.Error())
	}

	return nil
}
