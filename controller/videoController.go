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
	VideoList []service.Video `json:"video_list"`
	NextTime  int64           `json:"next_time"`
}

type VideoListResponse struct {
	Response
	VideoList []service.Video `json:"video_list"`
}

// Feed /feed/
func Feed(c *gin.Context) {
	inputTime := c.Query("latest_time")
	log.Printf("传入的时间" + inputTime)
	var lastTime time.Time
	if inputTime != "0" {
		me, _ := strconv.ParseInt(inputTime, 10, 64)
		lastTime = time.Unix(me, 0)
	} else {
		lastTime = time.Now()
	}
	log.Printf("获取到时间戳%v", lastTime)
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
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

// Publish /publish/action/
func Publish(c *gin.Context) {
	data, err := c.FormFile("data")
	userId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	title := c.PostForm("title")
	log.Printf("获取到视频title:%v\n", title)
	if err != nil {
		log.Printf("获取视频流失败:%v", err)
		c.JSON(http.StatusOK, Response{
			StatusCode: 1,
			StatusMsg:  err.Error(),
		})
		return
	}

	videoService := GetVideo()
	err = videoService.Publish(data, userId, title)
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
		StatusMsg:  "uploaded successfully",
	})
}

// PublishList /publish/list/
func PublishList(c *gin.Context) {
	user_Id, _ := c.GetQuery("user_id")
	userId, _ := strconv.ParseInt(user_Id, 10, 64)
	log.Printf("获取到用户id:%v\n", userId)
	curId, _ := strconv.ParseInt(c.GetString("userId"), 10, 64)
	log.Printf("获取到当前用户id:%v\n", curId)
	videoService := GetVideo()
	list, err := videoService.List(userId, curId)
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

// GetVideo 拼装videoService
func GetVideo() service.VideoServiceImpl {
	var userService service.UserServiceImpl
	var followService service.FollowServiceImp
	var videoService service.VideoServiceImpl
	var likeService service.LikeServiceImpl
	var commentService service.CommentServiceImpl
	userService.FollowService = &followService
	userService.LikeService = &likeService
	followService.UserService = &userService
	likeService.VideoService = &videoService
	commentService.UserService = &userService
	videoService.CommentService = &commentService
	videoService.LikeService = &likeService
	videoService.UserService = &userService
	return videoService
}
