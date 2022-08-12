package gojobs

import (
	"context"
	"github.com/robfig/cron/v3"
	"go.dtapp.net/dorm"
	"log"
	"testing"
)

func TestMasterRedis(t *testing.T) {

	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	// 创建一个cron实例 精确到秒
	c := cron.New(cron.WithSeconds())

	// 每隔5秒执行一次
	_, _ = c.AddFunc(GetSeconds(5).Spec(), func() {

		log.Println("每隔5秒执行一次")

		publish, err := client.Db.Publish(context.Background(), "test_cron_127.0.0.1", "每隔5秒执行一次").Result()
		t.Log(publish)
		t.Log(err)
	})

	// 每隔10秒执行一次
	_, _ = c.AddFunc(GetSeconds(10).Spec(), func() {

		log.Println("每隔10秒执行一次")

		publish, err := client.Db.Publish(context.Background(), "test_cron_127.0.0.1", "每隔10秒执行一次").Result()
		t.Log(publish)
		t.Log(err)
	})

	// 启动任务
	c.Start()

	// 关闭任务
	defer c.Stop()
	select {}
}
