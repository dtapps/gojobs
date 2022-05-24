package gojobs

import (
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
	client     *clientv3.Client // 驱动
	kv         clientv3.KV
	lease      clientv3.Lease
}

// Close 关闭
func (e Etcd) Close() {
	e.client.Close()
}
