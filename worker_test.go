package gojobs

import (
	"context"
	"go.dtapp.net/dorm"
	"testing"
)

func TestWorkerRedis(t *testing.T) {

	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	// 订阅channel1这个channel
	sub := client.Db.Subscribe(context.Background(), "test_cron_127.0.0.1")
	t.Log(sub)

	for msg := range sub.Channel() {

		// 打印收到的消息
		t.Log(msg)
		t.Log(msg.Channel)
		t.Log(msg.Payload)

		// 检测收到的消息类型
		//switch iface.(type) {
		//case *redis.Subscription:
		//	t.Log("订阅成功")
		//case *redis.Message:
		//	// 处理收到的消息
		//	// 这里需要做一下类型转换
		//	m := iface.(redis.Message)
		//	// 打印收到的小
		//	fmt.Println(m.Payload)
		//	t.Log("打印收到的小", m.Payload)
		//case *redis.Pong:
		//	t.Log("收到Pong消息")
		//default:
		//	// handle error
		//}
	}
}
