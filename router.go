package main

import (
	"TikTok/controller"
	"TikTok/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.GET("/feed/", middleware.AuthWithoutLogin(), controller.Feed)
	apiRouter.POST("/publish/action/", middleware.AuthBody(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.Auth(), controller.PublishList)
	apiRouter.GET("/user/", middleware.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.Auth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", middleware.Auth(), controller.GetFavouriteList)
	apiRouter.POST("/comment/action/", middleware.Auth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", middleware.Auth(), controller.CommentList)
	// extra apis - II
	apiRouter.POST("/relation/action/", middleware.Auth(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.Auth(), controller.GetFollowing)
	apiRouter.GET("/relation/follower/list", middleware.Auth(), controller.GetFollowers)
}
