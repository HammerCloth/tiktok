package main

import (
	"TikTok/dao"
	"TikTok/middleware/ffmpeg"
	"TikTok/middleware/ftp"
	"TikTok/middleware/rabbitmq"
	"TikTok/middleware/redis"
	"TikTok/util"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

//如果启动有问题，大概是你的IP地址出现变化，需要在项目依赖的服务器中配置安全组
func main() {
	//关闭log
	//log.SetOutput(ioutil.Discard)
	initDeps()
	//gin
	r := gin.Default()
	initRouter(r)
	//pprof
	pprof.Register(r)
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

// 加载项目依赖
func initDeps() {
	// 初始化数据库
	dao.Init()
	// 初始化FTP服务器链接
	ftp.InitFTP()
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
	//初始化Comment的消息队列，并开启消费
	rabbitmq.InitCommentRabbitMQ()
	//初始化敏感词拦截器。
	util.InitFilter()
}
