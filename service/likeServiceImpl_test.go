package service

import (
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"TikTok/middleware/rabbitmq"
	"TikTok/middleware/redis"
	"fmt"
	"testing"
)

func TestIsFavourite(t *testing.T) {
	// 初始化数据库
	dao.Init()
	// 初始化FTP服务器链接
	dao.InitFTP()
	// 初始化SSH
	ffmpeg.InitSSH()

	// 初始化redis-DB0的连接，follow选择的DB0.
	redis.InitRedis()
	// 初始化rabbitMQ。
	rabbitmq.InitRabbitMQ()
	// 初始化Follow的相关消息队列，并开启消费。
	rabbitmq.InitFollowRabbitMQ()
	// 初始化Like的相关消息队列，并开启消费。
	impl := LikeServiceImpl{}
	isFavourite, _ := impl.IsFavourite(666, 3)
	fmt.Printf("%v", isFavourite)
}

func TestFavouriteCount(t *testing.T) {
	// 初始化数据库
	dao.Init()
	// 初始化FTP服务器链接
	dao.InitFTP()
	// 初始化SSH
	ffmpeg.InitSSH()

	// 初始化redis-DB0的连接，follow选择的DB0.
	redis.InitRedis()
	// 初始化rabbitMQ。
	rabbitmq.InitRabbitMQ()
	// 初始化Follow的相关消息队列，并开启消费。
	rabbitmq.InitFollowRabbitMQ()
	// 初始化Like的相关消息队列，并开启消费。
	impl := LikeServiceImpl{}
	count, _ := impl.FavouriteCount(666)
	fmt.Printf("%v", count)
}

func TestTotalFavourite(t *testing.T) {
	// 初始化数据库
	dao.Init()
	// 初始化FTP服务器链接
	dao.InitFTP()
	// 初始化SSH
	ffmpeg.InitSSH()

	// 初始化redis-DB0的连接，follow选择的DB0.
	redis.InitRedis()
	// 初始化rabbitMQ。
	rabbitmq.InitRabbitMQ()
	// 初始化Follow的相关消息队列，并开启消费。
	rabbitmq.InitFollowRabbitMQ()
	// 初始化Like的相关消息队列，并开启消费。
	rabbitmq.InitLikeRabbitMQ()
	impl := LikeServiceImpl{}
	count, _ := impl.TotalFavourite(3)
	fmt.Printf("%v", count)
}
