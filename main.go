package main

import (
	"TikTok/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	dao.Init()
	// 初始化FTP服务器链接
	dao.InitFTP()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
