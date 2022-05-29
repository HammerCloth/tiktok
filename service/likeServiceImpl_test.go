package service

import (
	"TikTok/dao"
	"TikTok/middleware"
	"fmt"
	"testing"
)

func TestIsFavourite(t *testing.T) {
	impl := LikeServiceImpl{}
	isFavourite, _ := impl.IsFavourite(666, 3)
	fmt.Printf("%v", isFavourite)
}

func TestFavouriteCount(t *testing.T) {
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
	middleware.InitSSH()

	// 初始化redis-DB0的连接，follow选择的DB0.
	middleware.InitRedis()
	// 初始化rabbitMQ。
	middleware.InitRabbitMQ()
	// 初始化Follow的相关消息队列，并开启消费。
	middleware.InitFollowRabbitMQ()
	// 初始化Like的相关消息队列，并开启消费。
	middleware.InitLikeRabbitMQ()
	impl := LikeServiceImpl{}
	count, _ := impl.TotalFavourite(3)
	fmt.Printf("%v", count)
}
