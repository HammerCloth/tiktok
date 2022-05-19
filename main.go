package main

import (
	"TikTok/dao"
	"github.com/gin-gonic/gin"
)

func main() {
	// 初始化数据库
	dao.Init()

	r := gin.Default()

	initRouter(r)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")

}
