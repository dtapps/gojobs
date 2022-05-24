package gojobs

import (
	"context"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/v3"
	"log"
	"time"
)

// NewEtcdWorker 创建 etcd Worker
func NewEtcdWorker(config *EtcdConfig) (*Etcd, error) {

	var (
		e   = &Etcd{}
		err error
	)

	e.Endpoints = config.Endpoints
	e.DialTimeout = config.DialTimeout

	e.Client, err = clientv3.New(clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: e.DialTimeout,
	})
	if err != nil {
		return nil, errors.New("连接失败：" + err.Error())
	}

	// 得到KV和Lease的API子集
	e.Kv = clientv3.NewKV(e.Client)
	e.Lease = clientv3.NewLease(e.Client)

	// 注册
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

	for {
		// 注册路径
		regKey = JobWorkerDir + e.LocalIP

		cancelFunc = nil

		// 创建租约
		leaseGrantResp, err = e.Lease.Grant(context.TODO(), 10)
		log.Println("创建租约")
		if err != nil {
			goto RETRY
		}

		// 自动续租
		keepAliveChan, err = e.Lease.KeepAlive(context.TODO(), leaseGrantResp.ID)
		log.Println("自动续租")
		if err != nil {
			goto RETRY
		}

		cancelCtx, cancelFunc = context.WithCancel(context.TODO())

		// 注册到etcd
		_, err = e.Kv.Put(cancelCtx, regKey, "", clientv3.WithLease(leaseGrantResp.ID))
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
