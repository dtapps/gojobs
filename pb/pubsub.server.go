package pb

import (
	"context"
	"go.dtapp.net/gojobs/pubsub"
	"log"
	"strings"
	"time"
)

type PubSubServerService struct {
	pub *pubsub.Publisher
	UnimplementedPubSubServer
}

func NewPubSubServerService() *PubSubServerService {
	return &PubSubServerService{
		// 新建一个Publisher对象
		pub: pubsub.NewPublisher(time.Millisecond*100, 10),
	}
}

// Publish 实现发布方法
func (p *PubSubServerService) Publish(ctx context.Context, arg *String) (*String, error) {
	log.Printf("[服务中转]：%v\n", arg.GetValue())
	// 发布消息
	p.pub.Publish(arg.GetValue())
	return &String{Value: arg.GetValue()}, nil
}

// Subscribe 实现订阅方法
func (p *PubSubServerService) Subscribe(arg *String, stream PubSub_SubscribeServer) error {

	// SubscribeTopic 增加一个使用函数过滤器的订阅者
	// func(v interface{}) 定义函数过滤的规则
	// SubscribeTopic 返回一个chan interface{}

	log.Printf("[服务中转]收到任务：%v\n", arg.GetValue())

	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		// 接收数据是string，并且key是以arg为前缀的
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})

	log.Println("[服务中转]工作线：", ch)
	log.Println("[服务中转]工作线数量：", p.pub.Len())

	// 服务器遍历chan，并将其中信息发送给订阅客户端
	for v := range ch {
		err := stream.Send(&String{Value: v.(string)})
		if err != nil {
			log.Println("[服务中转]任务分配失败：", err.Error())
			return err
		}
	}

	return nil
}
