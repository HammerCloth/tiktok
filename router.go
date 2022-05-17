package main

import (
	"TikTok/controller"
	"TikTok/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")

	apiRouter.GET("/user/", middleware.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
}
