package main

import (
	"TikTok/dao"
	"TikTok/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

func main() {
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
	middleware.InitLikeRabbitMQ()

	//gin
	r := gin.Default()
	initRouter(r)

	//pprof
	pprof.Register(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
