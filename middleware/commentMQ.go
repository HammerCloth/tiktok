package middleware

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/streadway/amqp"
	"log"
	"strconv"
	"strings"
)

type CommentMQ struct {
	RabbitMQ
	channel   *amqp.Channel
	queueName string
	exchange  string
	key       string
}

var MqCommentAdd *CommentMQ //添加评论信息
var MqCommentDel *CommentMQ //删除评论信息

// InitCommentRabbitMQ 初始化rabbitMQ连接。
func InitCommentRabbitMQ() {
	MqCommentAdd = NewCommentMQ("comment_add")
	go MqCommentAdd.CommentConsumer()

	MqCommentDel = NewCommentMQ("comment_del")
	go MqCommentDel.CommentConsumer()
}

// NewCommentMQ 获取commentMQ的对应队列。
func NewCommentMQ(queueName string) *CommentMQ {
	commentMQ := &CommentMQ{
		RabbitMQ:  *Rmq,
		queueName: queueName,
	}

	cha, err := commentMQ.conn.Channel()
	commentMQ.channel = cha
	Rmq.failOnErr(err, "获取通道失败")
	return commentMQ
}

// CommentPublish
// 在对应的队列中将消息打包为byte，
func (c *CommentMQ) CommentPublish(body string) {
	_, err := c.channel.QueueDeclare(c.queueName, false, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	c.channel.Publish(c.exchange, c.queueName, false, false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})

}

// CommentConsumer Comment关系的消费逻辑。
func (c *CommentMQ) CommentConsumer() {
	log.Println("CommentConsumer running")
	//QueueDeclare声明一个队列来保存消息并传递给使用者。如果队列不存在，则声明将创建队列，或确保现有队列与相同的参数匹配.
	_, err := c.channel.QueueDeclare(c.queueName, false, false, false, false, nil)
	if err != nil {
		log.Println(err)
	}

	//2、接收消息
	msg, err := c.channel.Consume(c.queueName, "", true, false, false, false, nil)
	if err != nil {
		//panic(err)
		log.Println(err)
	}

	commentChan := make(chan bool)
	switch c.queueName {
	case "comment_add":
		go c.consumerCommentAdd(msg)
	case "comment_del":
		go c.consumerCommentDel(msg)
	}
	//log.Printf("[*] Waiting for messages,To exit press CTRL+C")
	<-commentChan
}

// 消费，评论添加-在redis中添加评论相关信息，共5个信息。
func (c *CommentMQ) consumerCommentAdd(messages <-chan amqp.Delivery) {
	log.Println("consumerCommentAdd running")
	for msg := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", msg.Body), " ")
		videoId, _ := strconv.Atoi(params[0])
		commentId, _ := strconv.Atoi(params[1])
		userId, _ := strconv.Atoi(params[2])
		text := params[3]
		date := params[4]

		//缓存评论id 使用zSet
		RdbVCid.ZAdd(Ctx, strconv.FormatInt(int64(videoId), 10), &redis.Z{
			Score:  float64(commentId),
			Member: commentId,
		})
		//缓存评论信息, comId: userId value, context value, date value。
		//hash
		RdbCInfo.HSet(Ctx, strconv.Itoa(commentId),
			"userId", strconv.Itoa(userId),
			"content", text, "date", date)
	}
}

// 消费，评论删除-在redis中删除评论相关信息，共2个信息，删除ZSet中的value和hash中的键值。
func (c *CommentMQ) consumerCommentDel(messages <-chan amqp.Delivery) {
	for msg := range messages {
		// 参数解析。
		params := strings.Split(fmt.Sprintf("%s", msg.Body), " ")
		videoId, _ := strconv.Atoi(params[0])
		commentId, _ := strconv.Atoi(params[1])

		//删除缓存中评论id
		_, err := RdbVCid.ZRem(Ctx, strconv.FormatInt(int64(videoId), 10)).Result()
		if err != nil {
			log.Println(err)
		}

		//删除缓存中评论信息, comId: userId value, context value, date value。
		//hash
		_, err = RdbCInfo.HDel(Ctx, strconv.Itoa(commentId)).Result()
		if err != nil {
			log.Println(err)
		}
	}
}

// 关闭mq通道。
func (c *CommentMQ) destroy() {
	c.channel.Close()
}
