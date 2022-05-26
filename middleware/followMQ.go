package middleware

import (
	"TikTok/dao"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
)

type FollowMQ struct {
	RabbitMQ
	queueName string
	exchange  string
	key       string
}

// NewFollowRabbitMQ 获取followMQ的对应队列。
func NewFollowRabbitMQ(queueName string) *FollowMQ {
	return &FollowMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}

}

// Publish follow关系的发布配置。
func (f *FollowMQ) Publish(message string) {

	_, err := f.channel.QueueDeclare(
		f.queueName,
		//是否持久化
		false,
		//是否为自动删除
		false,
		//是否具有排他性
		false,
		//是否阻塞
		false,
		//额外属性
		nil,
	)
	if err != nil {
		panic(err)
	}

	f.channel.Publish(
		f.exchange,
		f.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

}

// Consumer follow关系的消费逻辑。
func (f *FollowMQ) Consumer() {

	_, err := f.channel.QueueDeclare(f.queueName, false, false, false, false, nil)

	if err != nil {
		panic(err)
	}

	//2、接收消息
	msgs, err := f.channel.Consume(
		f.queueName,
		//用来区分多个消费者
		"",
		//是否自动应答
		true,
		//是否具有排他性
		false,
		//如果设置为true，表示不能将同一个connection中发送的消息传递给这个connection中的消费者
		false,
		//消息队列是否阻塞
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	forever := make(chan bool)
	switch f.queueName {
	case "follow_add":
		go f.consumerFollowAdd(msgs)
	case "follow_del":
		go f.consumerFollowDel(msgs)

	}

	log.Printf("[*] Waiting for messagees,To exit press CTRL+C")

	<-forever

}

// 关系添加的消费方式。
func (f *FollowMQ) consumerFollowAdd(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		targetId, _ := strconv.Atoi(params[1])
		// 日志记录。
		sql := fmt.Sprintf("CALL addFollowRelation(%v,%v)", targetId, userId)
		log.Printf("消费队列执行添加关系。SQL如下：%s", sql)
		// 执行SQL，注必须scan，该SQL才能被执行。
		if err := dao.Db.Raw(sql).Scan(nil).Error; nil != err {
			// 执行出错，打印日志。
			log.Println(err.Error())
		}
	}
}

// 关系删除的消费方式。
func (f *FollowMQ) consumerFollowDel(msgs <-chan amqp.Delivery) {
	for d := range msgs {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", d.Body), " ")
		userId, _ := strconv.Atoi(params[0])
		targetId, _ := strconv.Atoi(params[1])
		// 日志记录。
		sql := fmt.Sprintf("CALL delFollowRelation(%v,%v)", targetId, userId)
		//log.Printf("消费队列执行删除关系。SQL如下：%s", sql)
		// 执行SQL，注必须scan，该SQL才能被执行。
		if err := dao.Db.Raw(sql).Scan(nil).Error; nil != err {
			// 执行出错，打印日志。
			log.Println(err.Error())
		}
	}
}

var RmqFollowAdd *FollowMQ
var RmqFollowDel *FollowMQ

// InitFollowRabbitMQ 初始化rabbitMQ连接。
func InitFollowRabbitMQ() {
	RmqFollowAdd = NewFollowRabbitMQ("follow_add")
	go RmqFollowAdd.Consumer()

	RmqFollowDel = NewFollowRabbitMQ("follow_del")
	go RmqFollowDel.Consumer()
}
