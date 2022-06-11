package gojobs

import (
	"context"
	"github.com/robfig/cron/v3"
	"log"
	"testing"
	"time"
)

func TestEtcdServer(t *testing.T) {
	server, err := NewEtcdServer(&EtcdConfig{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
		Username:    "root",
		Password:    "p5sttPYcFWw7Z7aP",
	})
	if err != nil {
		panic(err)
	}
	defer server.Close()

	// 创建一个cron实例 精确到秒
	c := cron.New(cron.WithSeconds())

	// 每隔15秒执行一次
	_, _ = c.AddFunc("*/15 * * * * *", func() {

		create, err := server.Create(context.TODO(), server.IssueWatchKey("116.30.228.12")+"/"+"wechat_1_test", "每隔15秒执行一次")
		if err != nil {
			log.Println("创建任务失败", err)
		}
		log.Println("创建任务成功", create, err)

	})

	// 每隔30秒执行一次
	_, _ = c.AddFunc("*/30 * * * * *", func() {

		create, err := server.Create(context.TODO(), server.IssueWatchKey("127.0.0.1")+"/"+"wechat_2_test", "每隔30秒执行一次")
		if err != nil {
			log.Println("创建任务失败", err)
		}
		log.Println("创建任务成功", create, err)

	})

	// 每隔1分钟执行一次
	_, _ = c.AddFunc("0 */1 * * * *", func() {

		create, err := server.Create(context.TODO(), server.IssueWatchKey("116.30.228.12")+"/"+"wechat_3_test", "每隔1分钟执行一次")
		if err != nil {
			log.Println("创建任务失败", err)
		}
		log.Println("创建任务成功", create, err)

		workers, _ := server.ListWorkers()
		log.Println("ListWorkers", workers)
	})

	// 启动任务
	c.Start()

	// 关闭任务
	defer c.Stop()
	select {}
}
