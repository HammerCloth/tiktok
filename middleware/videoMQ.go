package middleware

import (
	"github.com/streadway/amqp"
	"log"
	"strings"
)

/**
这里创建队列，用于调用ffmpeg防止因为并发过多而出现问题
*/

var VideoChannel *amqp.Channel

const nameEX = "videoEX"
const nameQueue = "ffmpeg"
const key = "ffmpeg"

func InitVideoMQ() error {
	var err error
	VideoChannel, err = Rmq.conn.Channel()
	if err != nil {
		log.Fatal("Failed to open a channel")
		return err
	}
	log.Println("open a channel Success")

	//exchange 交换器的名称
	//type 交换器的类型，常见的有fanout、direct、topic、headers这四种
	//durable 设置是否持久 durab 设置为 true 表示持久化， 反之是非持久,设置为true则将Exchange存盘，即使服务器重启数据也不会丢失
	//autoDelete 设置是否自动删除，当最后一个绑定到Exchange上的队列删除后，自动删除该Exchange，简单来说也就是如果该Exchange没有和任何队列Queue绑定则删除
	//internal 设置是否是RabbitMQ内部使用，默认false。如果设置为 true ，则表示是内置的交换器，客户端程序无法直接发送消息到这个交换器中，只能通过交换器路由到交换器这种方式。
	//argument 扩展参数，用于扩展AMQP协议自制定化使用
	//noWait:是否非阻塞，true表示是。阻塞：表示创建交换器的请求发送后，阻塞等待RMQ Server返回信息。
	//非阻塞：不会阻塞等待RMQ Server的返回信息，而RMQ Server也不会返回信息。（不推荐使用）
	err = VideoChannel.ExchangeDeclarePassive(nameEX, amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare Exchange")
		return err
	}
	log.Println("declare Exchange Success")
	//queue ：队列的名称
	//durable ：设置是否持久化。为true 则设置队列为持久化。持久化的队列会存盘，在服务器重启的时候可以保证不丢失相关信息
	//exclusive ：设置是否排他（独占队列）。为true 则设置队列为排他的。如果一个队列被声明为排他队列，该队列仅对首次声明它的连接可见，并在连接断开时自动删除。这里需要注意三点:排他队列是基于连接( Connection) 可见的，同一个连接的不同信道(Channel)是可以同时访问同一连接创建的排他队列; "首次"是指如果一个连接己经声明了一个排他队列，其他连接是不允许建立同名的排他队列的，这个与普通队列不同:即使该队列是持久化的，一旦连接关闭或者客户端退出，该排他队列都会被自动删除，这种队列适用于一个客户端同时发送和读取消息的应用场景
	//autoDelete ：设置是否自动删除，当最后一个监听被移除后，自动删除队列；也就是说至少有一个消费者连接到这个队列，之后所有与这个队列连接的消费者都断开时，才会自动删除
	//arguments ：设置队列的一些其他参数；
	_, err = VideoChannel.QueueDeclare(nameQueue, true, false, false, false, nil)
	if err != nil {
		log.Fatal("Failed to declare Queue")
		return err
	}
	log.Println("declare Queue Success")
	//绑定
	err = VideoChannel.QueueBind(nameQueue, key, nameEX, false, nil)
	if err != nil {
		log.Fatal("Failed to banding")
		return err
	}
	log.Println("binding Success")
	return nil
}

func MQffmpeg(videoName string, imageName string) error {
	body := videoName + "," + imageName
	err := VideoChannel.Publish(nameEX, key, false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
	if err != nil {
		log.Fatalf("消息%v发送失败", body)
		return err
	}
	log.Printf("消息%v发送成功", body)
	return nil
}

func Exeffmpeg() {
	msgs, _ := VideoChannel.Consume(nameQueue, // queue
		"",    // consumer
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			s := string(d.Body)
			names := strings.Split(s, ",")
			err := Ffmpeg(names[0], names[1])
			if err != nil {
				return
			}
		}
	}()
}
