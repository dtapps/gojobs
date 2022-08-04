package gojobs

import (
	"context"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"strings"
)

// NewEtcdServer 创建 etcd server
func NewEtcdServer(config *EtcdConfig) (*Etcd, error) {

	var (
		e   = &Etcd{}
		err error
	)

	e.Endpoints = config.Endpoints
	e.DialTimeout = config.DialTimeout
	e.LocalIP = config.LocalIP
	e.Username = config.Username
	e.Password = config.Password
	e.CustomDirectory = config.CustomDirectory

	e.Debug = config.Debug

	v3Config := clientv3.Config{
		Endpoints:   e.Endpoints,
		DialTimeout: e.DialTimeout,
	}

	// 判断有没有配置用户信息
	if e.Username != "" {
		v3Config.Username = e.Username
		v3Config.Password = e.Password
	}

	e.Client, err = clientv3.New(v3Config)
	if err != nil {
		return nil, errors.New("连接失败：" + err.Error())
	}

	// kv API子集
	e.Kv = clientv3.NewKV(e.Client)

	// 创建一个lease（租约）对象
	e.Lease = clientv3.NewLease(e.Client)

	return e, nil
}

// ListWorkers 获取在线worker列表
func (e Etcd) ListWorkers() (workerArr []string, err error) {
	var (
		getResp  *clientv3.GetResponse
		kv       *mvccpb.KeyValue
		workerIP string
	)

	// 初始化数组
	workerArr = make([]string, 0)

	// 获取目录下所有Kv
	getResp, err = e.Kv.Get(context.TODO(), getJobWorkerDir(e), clientv3.WithPrefix())
	if err != nil {
		return workerArr, err
	}

	// 解析每个节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key : /cron/workers/192.168.2.1
		workerIP = e.ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}
	return workerArr, err
}

// ExtractWorkerIP 提取worker的IP
func (e Etcd) ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, getJobWorkerDir(e))
}
