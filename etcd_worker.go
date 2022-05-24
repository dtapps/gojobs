package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.dtapp.net/goip"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

const (
	// JobWorkerDir 服务注册目录
	JobWorkerDir = "/cron/workers/"
)

// NewEtcdWorker 创建 etcd Worker
func NewEtcdWorker(config *EtcdConfig) (e *Etcd, err error) {

	e.Endpoints = config.Endpoints
	e.DialTimeout = config.DialTimeout

	e.client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: e.DialTimeout,
	})
	if err != nil {
		return nil, errors.New("连接失败：" + err.Error())
	}

	// 得到KV和Lease的API子集
	e.kv = clientv3.NewKV(e.client)
	e.lease = clientv3.NewLease(e.client)

	go e.RegisterWorker()

	return e, nil
}

// RegisterWorker 注册worker
func (e Etcd) RegisterWorker() {
	var (
		regKey         string
		leaseGrantResp *clientv3.LeaseGrantResponse
		err            error
		keepAliveChan  <-chan *clientv3.LeaseKeepAliveResponse
		keepAliveResp  *clientv3.LeaseKeepAliveResponse
		cancelCtx      context.Context
		cancelFunc     context.CancelFunc
	)

	localIP := goip.GetOutsideIp()

	for {
		// 注册路径
		regKey = JobWorkerDir + localIP

		cancelFunc = nil

		// 创建租约
		leaseGrantResp, err = e.lease.Grant(context.TODO(), 10)
		if err != nil {
			log.Println("创建租约")
			goto RETRY
		}

		// 自动续租
		keepAliveChan, err = e.lease.KeepAlive(context.TODO(), leaseGrantResp.ID)
		if err != nil {
			log.Println("自动续租")
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		// 注册到etcd
		_, err = e.kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID))
		if err != nil {
			log.Println(fmt.Sprintf(" %s 服务注册失败:%s", regKey, err))
			goto RETRY
		}

		// 处理续租应答
		for {
			select {
			case keepAliveResp = <-keepAliveChan:
				if keepAliveResp == nil { // 续租失败
					log.Println("续租失败")
					goto RETRY
				}
			}
		}

	RETRY:
		time.Sleep(1 * time.Second)
		if cancelFunc != nil {
			cancelFunc()
		}
	}
}
