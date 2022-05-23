package gojobs

import (
	"context"
	"go.dtapp.net/gojobs/pb"
	"io"
	"log"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestJobs(t *testing.T) {
	wg := sync.WaitGroup{}
	wg.Add(3)
	go testServer(&wg)
	go testCron(&wg)
	go testWorker(&wg)
	wg.Wait()
}

func testServer(wg *sync.WaitGroup) {

	server := NewServer(&ServerConfig{
		PublishTimeout: time.Millisecond * 100,
		PubBuffer:      10,
		Address:        "0.0.0.0:8888",
	})

	cronServer := server.Pub.SubscribeTopic(func(v interface{}) bool {
		log.Println("SubscribeTopic:", v)
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, "cron:") {
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

	t1 := time.NewTimer(time.Second * 10)
	for {
		select {
		case <-t1.C:

			server.Send(&pb.PublishRequest{
				Value: "cron:" + "wechat.send",
			})

			t1.Reset(time.Second * 10)
		}
	}

	wg.Done()
}

func testWorker(wg *sync.WaitGroup) {

	server := NewCron(&CronConfig{
		Address: "localhost:8888",
	})
	defer server.Conn.Close()

	// 订阅服务，传入参数是 cron:
	// 会想过滤器函数，订阅者应该收到的信息为 cron:任务名称
	stream, err := server.Pub.Subscribe(context.Background(), &pb.SubscribeRequest{
		Value: "cron:",
		Ip:    "127.0.0.1",
	})
	if err != nil {
		log.Printf("[跑业务]发送失败:%v\n", err)
	}

	// 阻塞遍历流，输出结果
	for {
		reply, err := stream.Recv()
		if io.EOF == err {
			log.Println("[跑业务]已关闭:", err.Error())
			break
		}
		if nil != err {
			log.Println("[跑业务]异常:", err.Error())
			break
		}
		log.Println("[跑业务]:", reply)
	}

	wg.Done()
}
