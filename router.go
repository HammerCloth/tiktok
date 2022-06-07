package main

import (
	"TikTok/controller"
	"TikTok/middleware/jwt"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	apiRouter := r.Group("/douyin")
	// basic apis
	apiRouter.GET("/feed/", jwt.AuthWithoutLogin(), controller.Feed)
	apiRouter.POST("/publish/action/", jwt.AuthBody(), controller.Publish)
	apiRouter.GET("/publish/list/", jwt.Auth(), controller.PublishList)
	apiRouter.GET("/user/", jwt.Auth(), controller.UserInfo)
	apiRouter.POST("/user/register/", controller.Register)
	apiRouter.POST("/user/login/", controller.Login)
	// extra apis - I
	apiRouter.POST("/favorite/action/", jwt.Auth(), controller.FavoriteAction)
	apiRouter.GET("/favorite/list/", jwt.Auth(), controller.GetFavouriteList)
	apiRouter.POST("/comment/action/", jwt.Auth(), controller.CommentAction)
	apiRouter.GET("/comment/list/", jwt.AuthWithoutLogin(), controller.CommentList)
	// extra apis - II
	apiRouter.POST("/relation/action/", jwt.Auth(), controller.RelationAction)
	apiRouter.GET("/relation/follow/list/", jwt.Auth(), controller.GetFollowing)
	apiRouter.GET("/relation/follower/list", jwt.Auth(), controller.GetFollowers)
}
