package middleware

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
)

const MQURL = "amqp://tiktok:tiktok@106.14.75.229:5672/"

type RabbitMQ struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	mqurl   string
}

var Rmq *RabbitMQ

// InitRabbitMQ 初始化RabbitMQ的连接和通道。
func InitRabbitMQ() {

	Rmq = &RabbitMQ{
		mqurl: MQURL,
	}
	dial, err := amqp.Dial(Rmq.mqurl)
	Rmq.failOnErr(err, "创建连接失败")
	Rmq.conn = dial

	Rmq.channel, err = Rmq.conn.Channel()
	Rmq.failOnErr(err, "获取通道失败")
}

// 连接出错时，输出错误信息。
func (r *RabbitMQ) failOnErr(err error, message string) {
	if err != nil {
		log.Fatalf("%s:%s\n", err, message)
		panic(fmt.Sprintf("%s:%s\n", err, message))
	}
}

// 关闭mq通道和mq的连接。
func (r *RabbitMQ) destory() {
	r.channel.Close()
	r.conn.Close()
}