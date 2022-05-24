package gojobs

import (
	"context"
	"errors"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"strings"
	"time"
)

// EtcdConfig etcd配置
type EtcdConfig struct {
	Endpoints   []string      // 接口 []string{"http://127.0.0.1:2379"}
	DialTimeout time.Duration // time.Second * 5
}

// Etcd etcd
type Etcd struct {
	EtcdConfig                  // 配置
	client     *clientv3.Client // 驱动
	kv         clientv3.KV
	lease      clientv3.Lease
}

// Close 关闭
func (e Etcd) Close() {
	e.client.Close()
}

// NewEtcd 创建etcd
func NewEtcd(config *EtcdConfig) (*Etcd, error) {

	e := &Etcd{}
	e.Endpoints = config.Endpoints
	e.DialTimeout = config.DialTimeout

	var err error
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

	return e, nil
}

// Watch 监听
func (e Etcd) Watch(ctx context.Context, key string, opts ...clientv3.OpOption) clientv3.WatchChan {
	return e.client.Watch(ctx, key, opts...)
}

// Create 创建
func (e Etcd) Create(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return e.client.Put(ctx, key, val, opts...)
}

// Get 获取
func (e Etcd) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return e.client.Get(ctx, key, opts...)
}

// Update 更新
func (e Etcd) Update(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return e.client.Put(ctx, key, val, opts...)
}

// Delete 删除
func (e Etcd) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return e.client.Delete(ctx, key, opts...)
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
	if getResp, err = e.kv.Get(context.TODO(), JobWorkerDir, clientv3.WithPrefix()); err != nil {
		return
	}

	// 解析每个节点的IP
	for _, kv = range getResp.Kvs {
		// kv.Key : /cron/workers/192.168.2.1
		workerIP = ExtractWorkerIP(string(kv.Key))
		workerArr = append(workerArr, workerIP)
	}
	return
}

const (
	// JobWorkerDir 服务注册目录
	JobWorkerDir = "/cron/workers/"
)

// ExtractWorkerIP 提取worker的IP
func ExtractWorkerIP(regKey string) string {
	return strings.TrimPrefix(regKey, JobWorkerDir)
}
