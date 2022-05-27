package main

import (
	"TikTok/controller"
	"TikTok/middleware"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.GET("/feed/", middleware.Auth(), controller.Feed)
	apiRouter.POST("/publish/action/", middleware.Auth_body(), controller.Publish)
	apiRouter.GET("/publish/list/", middleware.Auth(), controller.PublishList)

	apiRouter.GET("/user/", middleware.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	// extra apis - I
	apiRouter.POST("/favorite/action/", middleware.Auth(), controller.Favorite_Action)
	apiRouter.GET("/favorite/list/", middleware.Auth(), controller.GetFavouriteList)
	//apiRouter.POST("/comment/action/", controller.CommentAction)
	//apiRouter.GET("/comment/list/", controller.CommentList)

	// extra apis - II
	//apiRouter.POST("/relation/action/", controller.RelationAction)
	//apiRouter.GET("/relation/follow/list/", controller.FollowList)
	//apiRouter.GET("/relation/follower/list/", controller.FollowerList)

	/*
		关注模块
	*/
	apiRouter.POST("/relation/action/", middleware.Auth(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", middleware.Auth(), controller.GetFollowing)
	apiRouter.GET("/relation/follower/list", middleware.Auth(), controller.GetFollowers)

	/*
		评论模块
	*/
	//发表评论
	apiRouter.POST("/comment/action/", middleware.Auth(), controller.CommentAction)
	//查看评论列表
	apiRouter.GET("/comment/list/", middleware.Auth(), controller.CommentList)
}
