package gojobs

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.dtapp.net/gojobs/pb"
	"go.dtapp.net/gouuid"
	"io"
	"log"
	"strings"
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

	cronServer := server.Pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, prefix) {
				return true
			}
		}
		return false
	})

	go func() {
		log.Println("cronServer：topic:", <-cronServer)
	}()

	err := server.StartUp()
	if err != nil {
		log.Panicf("创建服务失败:%v\n", err)
		return
	}

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
			Id:    gouuid.GetUuId(),
			Value: prefix + "wechat.send" + " 我是定时任务",
			Ip:    "127.0.0.1",
		})

	})

	// 每隔30秒执行一次
	_, _ = c.AddFunc("*/30 * * * * *", func() {

		server.Send(&pb.PublishRequest{
			Id:    gouuid.GetUuId(),
			Value: prefix + "wechat.send" + " 我是定时任务",
			Ip:    "14.155.157.19",
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

	server := NewCron(&CronConfig{
		Address: "localhost:8888",
	})
	defer server.Conn.Close()

	// 订阅服务，传入参数是 cron:
	// 会想过滤器函数，订阅者应该收到的信息为 cron:任务名称
	stream, err := server.Pub.Subscribe(context.Background(), &pb.SubscribeRequest{
		Id:    gouuid.GetUuId(),
		Value: prefix,
		Ip:    "127.0.0.1",
	})
	if err != nil {
		log.Printf("[跑业务1]发送失败:%v\n", err)
	}

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
		log.Println("[跑业务1]:", reply)
	}

	wg.Done()
}

func testWorker2(wg *sync.WaitGroup) {

	server := NewCron(&CronConfig{
		Address: "localhost:8888",
	})
	defer server.Conn.Close()

	// 订阅服务，传入参数是 cron:
	// 会想过滤器函数，订阅者应该收到的信息为 cron:任务名称
	stream, err := server.Pub.Subscribe(context.Background(), &pb.SubscribeRequest{
		Id:    gouuid.GetUuId(),
		Value: prefix,
		Ip:    "14.155.157.19",
	})
	if err != nil {
		log.Printf("[跑业务2]发送失败:%v\n", err)
	}

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
		log.Println("[跑业务2]:", reply)
	}

	wg.Done()
}
