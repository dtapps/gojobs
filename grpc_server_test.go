package gojobs

import (
	"github.com/robfig/cron/v3"
	"go.dtapp.net/gojobs/pb"
	"go.dtapp.net/gouuid"
	"io"
	"log"
	"sync"
	"testing"
	"time"
)

func TestJobs(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(4)
	go testServer(&wg)
	go testCron(&wg)
	go testWorker1(&wg)
	go testWorker2(&wg)
	wg.Wait()
}

func testServer(wg *sync.WaitGroup) {

	server := NewServer(&ServerConfig{
		PublishTimeout: time.Millisecond * 100,
		PubBuffer:      10,
		Address:        "0.0.0.0:8888",
	})

	// 启动定时任务
	server.StartCron()

	// 启动服务
	server.StartUp()

	<-make(chan bool)

	wg.Done()
}

func testCron(wg *sync.WaitGroup) {

	server := NewCron(&CronConfig{
		Address: "localhost:8888",
	})
	defer server.Conn.Close()

	// 创建一个cron实例 精确到秒
	c := cron.New(cron.WithSeconds())

	// 每隔15秒执行一次
	_, _ = c.AddFunc("*/15 * * * * *", func() {

		server.Send(&pb.PublishRequest{
			Id:     gouuid.GetUuId(),
			Value:  prefix,
			Method: "wechat.1.send",
			Ip:     "127.0.0.1",
		})

	})

	// 每隔30秒执行一次
	_, _ = c.AddFunc("*/30 * * * * *", func() {

		server.Send(&pb.PublishRequest{
			Id:     gouuid.GetUuId(),
			Value:  prefix,
			Method: "wechat.2.send",
			Ip:     "14.155.157.19",
		})

	})

	// 启动任务
	c.Start()

	// 关闭任务
	defer c.Stop()
	select {}

	wg.Done()
}

func testWorker1(wg *sync.WaitGroup) {

	server := NewWorker(&WorkerConfig{
		Address:  "localhost:8888",
		ClientIp: "127.0.0.1",
	})
	defer server.Conn.Close()

	// 订阅服务
	stream := server.SubscribeCron()

	// 启动任务，会想过滤器函数，订阅者应该收到的信息为 cron:任务名称
	//stream := server.StartCron()

	// 阻塞遍历流，输出结果
	for {
		reply, err := stream.Recv()
		if io.EOF == err {
			log.Println("[跑业务1]已关闭:", err.Error())
			break
		}
		if nil != err {
			log.Println("[跑业务1]异常:", err.Error())
			break
		}
		log.Printf("[跑业务1]{收到}编号：%s 方法：%s\n", reply.GetId(), reply.GetMethod())
	}

	wg.Done()
}

func testWorker2(wg *sync.WaitGroup) {

	server := NewWorker(&WorkerConfig{
		Address:  "localhost:8888",
		ClientIp: "14.155.157.19",
	})
	defer server.Conn.Close()

	// 订阅服务
	stream := server.SubscribeCron()

	// 启动任务，会想过滤器函数，订阅者应该收到的信息为 cron:任务名称
	//stream := server.StartCron()

	// 阻塞遍历流，输出结果
	for {
		reply, err := stream.Recv()
		if io.EOF == err {
			log.Println("[跑业务2]已关闭:", err.Error())
			break
		}
		if nil != err {
			log.Println("[跑业务2]异常:", err.Error())
			break
		}
		log.Printf("[跑业务2]{收到}编号：%s 方法：%s\n", reply.GetId(), reply.GetMethod())
	}

	wg.Done()
}
