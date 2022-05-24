package gojobs

import (
	"context"
	"go.etcd.io/etcd/client/v3"
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
	c          *clientv3.Client // 驱动
}

// Close 关闭
func (e Etcd) Close() {
	e.c.Close()
}

// NewEtcd 创建etcd
func NewEtcd(config *EtcdConfig) *Etcd {
	e := &Etcd{}

	var err error
	e.c, err = clientv3.New(clientv3.Config{
		Endpoints:   config.Endpoints,
		DialTimeout: config.DialTimeout,
	})
	if err != nil {
		panic("连接失败：" + err.Error())
	}

	return e
}

// Watch 监听
func (e Etcd) Watch(ctx context.Context, key string) clientv3.WatchChan {
	return e.c.Watch(ctx, key) // type WatchChan <-chan WatchResponse
}

// Create 创建
func (e Etcd) Create(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return e.c.Put(ctx, key, val, opts...)
}

// Get 获取
func (e Etcd) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	return e.c.Get(ctx, key, opts...)
}

// Update 更新
func (e Etcd) Update(ctx context.Context, key, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	return e.c.Put(ctx, key, val, opts...)
}

// Delete 删除
func (e Etcd) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	return e.c.Delete(ctx, key, opts...)
}
