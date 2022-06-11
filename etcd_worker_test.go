package gojobs

import (
	"context"
	"fmt"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"log"
	"sync"
	"testing"
	"time"
)

func TestEtcdWorker(t *testing.T) {
	server, err := NewEtcdWorker(&EtcdConfig{
		Endpoints:   []string{"http://127.0.0.1:2379"},
		DialTimeout: time.Second * 5,
		Username:    "root",
		Password:    "p5sttPYcFWw7Z7aP",
	})
	if err != nil {
		panic(err)
	}
	defer server.Close()

	// 监听
	go func() {
		rch := server.Watch(context.TODO(), server.GetWatchKey()+"/", clientv3.WithPrefix())
		// 处理监听事件
		for watchResp := range rch {
			for _, watchEvent := range watchResp.Events {
				switch watchEvent.Type {
				case mvccpb.PUT:

					// 收到任务
					log.Printf("监听收到任务 键名:%s 值:%s\n", watchEvent.Kv.Key, watchEvent.Kv.Value)
					log.Println("处理监听")

					wg := sync.WaitGroup{}
					wg.Add(1)
					go run(&wg)
					wg.Wait()

					_, err = server.Delete(context.TODO(), string(watchEvent.Kv.Key))
					if err != nil {
						log.Println("删除失败", err)
					} else {
						log.Println("删除成功")
					}
				case mvccpb.DELETE:
					// 任务被删除了

				}

			}
		}
		log.Println("out")
	}()

	select {}
}

func run(wg *sync.WaitGroup) {
	log.Println("等待开始")
	fmt.Println(time.Now())
	time.Sleep(time.Second * 10)
	fmt.Println(time.Now())
	log.Println("等待结束")
	wg.Done()
}
