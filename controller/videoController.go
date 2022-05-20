package controller

import (
	"TikTok/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"time"
)

type FeedResponse struct {
	Response
	VideoList []service.Video `json:"video_list,omitempty"`
	NextTime  int64           `json:"next_time,omitempty"`
}

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

func GetVideo() service.VideoServiceImpl {
	var userService service.UserServiceImpl
	var followService service.FollowServiceImp
	var videoService service.VideoServiceImpl
	var likeService service.LikeServiceImpl
	var commentService service.CommentServiceImpl
	userService.FollowService = &followService
	followService.UserService = &userService
	likeService.VideoService = &videoService
	commentService.UserService = &userService
	videoService.CommentService = &commentService
	videoService.LikeService = &likeService
	videoService.UserService = &userService
	return videoService
}

func Feed(c *gin.Context) {
	last_time, _ := strconv.ParseInt(c.Query("latest_time"), 10, 64)
	lastTime := time.Unix(last_time, 0)
	log.Printf("获取到时间戳%v", lastTime)
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	videoService := GetVideo()
	feed, nextTime, err := videoService.Feed(lastTime, userId)
	if err != nil {
		log.Printf("方法videoService.Feed(lastTime, userId) 失败：%v", err)
		c.JSON(http.StatusOK, FeedResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频流失败"},
		})
		return
	}
	log.Printf("方法videoService.Feed(lastTime, userId) 成功")
	c.JSON(http.StatusOK, FeedResponse{
		Response:  Response{StatusCode: 0},
		VideoList: feed,
		NextTime:  nextTime.Unix(),
	})
}

// Publish apiRouter.POST("/publish/action/", controller.Publish)
func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	if err != nil {
		log.Printf("获取视频流失败:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	videoService := GetVideo()
	err = videoService.Publish(data, userId)
	if err != nil {
		log.Printf("方法videoService.Publish(data, userId) 失败：%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}
	log.Printf("方法videoService.Publish(data, userId) 成功")
	c.JSON(http.StatusOK, Response{
		StatusCode: 0,
		StatusMsg:  " uploaded successfully",
	})
}

// PublishList apiRouter.GET("/publish/list/", controller.PublishList)
func PublishList(c *gin.Context) {
	userId, _ := strconv.ParseInt(c.Query("user_id"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	videoService := GetVideo()
	list, err := videoService.List(userId)
	if err != nil {
		log.Printf("调用videoService.List(%v)出现错误：%v\n", userId, err)
		c.JSON(http.StatusOK, VideoListResponse{
			Response: Response{StatusCode: 1, StatusMsg: "获取视频列表失败"},
		})
		return
	}
	log.Printf("调用videoService.List(%v)成功", userId)
	c.JSON(http.StatusOK, VideoListResponse{
		Response:  Response{StatusCode: 0},
		VideoList: list,
	})
}
