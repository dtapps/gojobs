package gojobs

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"go.dtapp.net/gotime"
	"testing"
)
import "go.dtapp.net/dorm"

// 下发
func TestRedisPublish(t *testing.T) {
	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	publish, err := client.Db.Publish(context.Background(), "ch_user_1", "测试").Result()
	t.Log(publish)
	t.Log(err)
}

// 订阅
func TestRedisPSubscribe(t *testing.T) {
	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	// 订阅channel1这个channel
	sub := client.Db.PSubscribe(context.Background(), "ch_user_*")
	t.Log(sub)
	// sub.Channel() 返回go channel，可以循环读取redis服务器发过来的消息
	for msg := range sub.Channel() {
		// 打印收到的消息
		t.Log(msg.Channel)
		t.Log(msg.Payload)
	}
}

// 添加一个或者多个元素到集合，如果元素已经存在则更新分数
func TestRedisZAdd(t *testing.T) {
	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	// 删除小于当前时间的客户端
	client.Db.ZRemRangeByScore(context.Background(), "cron_jobs", "0", fmt.Sprintf("(%v", gotime.Current().Timestamp())).Result()

	// 添加客户端或者更新客户端
	publish, err := client.Db.ZAdd(context.Background(), "cron_jobs", redis.Z{
		Score:  float64(gotime.Current().AfterMinute(10).Timestamp()),
		Member: "127.0.0.1",
	}).Result()
	t.Log(publish)
	t.Log(err)
}

// 返回集合中某个索引范围的元素，根据分数从小到大排序
func TestRedisZRange(t *testing.T) {
	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	vals, err := client.Db.ZRange(context.Background(), "cron_jobs", 0, -1).Result()
	t.Log(err)

	for _, val := range vals {
		t.Log(val)
	}

}

func TestRedisKeys(t *testing.T) {
	client, err := dorm.NewRedisClient(&dorm.ConfigRedisClient{
		Addr:     "119.29.14.159:6379",
		Password: "980202",
		DB:       5,
	})
	if err != nil {
		t.Error(err)
	}

	client.Db.Set(context.Background(), "cron_jobs_client_127.0.0.1", "127.0.0.1", 0)
	client.Db.Set(context.Background(), "cron_jobs_client_43.135.79.235", "43.135.79.235", 0)

	result, err := client.Keys(context.Background(), "cron_jobs_client_*").Result()
	t.Log(result)
	t.Log(err)
	for _, val := range result {
		t.Log(val)
	}
}
