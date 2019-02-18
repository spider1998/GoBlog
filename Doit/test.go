package main
import (
	"fmt"
	"time"
	"github.com/nsqio/go-nsq"
)
// nsq发布消息
func Producer() {
	p, err := nsq.NewProducer("127.0.0.1:4150", nsq.NewConfig())                // 新建生产者
	if err != nil {
		panic(err)
	}
	if err := p.Publish("test", []byte("hello NSQ!!!")); err != nil {           // 发布消息
		panic(err)
	}
}
// nsq订阅消息
type ConsumerT struct{}
func (*ConsumerT) HandleMessage(msg *nsq.Message) error {
	fmt.Println(string(msg.Body))
	return nil
}
func Consumer() {
	c, err := nsq.NewConsumer("test", "test-channel", nsq.NewConfig())   // 新建一个消费者
	if err != nil {
		panic(err)
	}
	c.AddHandler(&ConsumerT{})                                           // 添加消息处理
	if err := c.ConnectToNSQD("127.0.0.1:4150"); err != nil {            // 建立连接
		panic(err)
	}
}
// 主函数
func main() {
	Producer()
	Consumer()
	time.Sleep(time.Second * 3)
}
// 运行将会打印： hello NSQ!!!