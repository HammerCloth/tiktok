package controller

import (
	"TikTok/service"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

type likeResponse struct {
	StatusCode int32  `json:"status_code"`
	StatusMsg  string `json:"status_msg,omitempty"`
}

type GetFavouriteListResponse struct {
	StatusCode int32           `json:"status_code"`
	StatusMsg  string          `json:"status_msg,omitempty"`
	Videolist  []service.Video `json:"video_list,omitempty"`
}

//点赞或者取消赞操作;
func Favorite_Action(c *gin.Context) {
	user_id := c.Query("user_id")
	userid, _ := strconv.ParseInt(user_id, 10, 64)
	video_id := c.Query("video_id")
	videoid, _ := strconv.ParseInt(video_id, 10, 64)
	action_type := c.Query("action_type")
	actiontype, _ := strconv.ParseInt(action_type, 10, 64)
	like := new(service.LikeServiceImpl)
	//获取点赞或者取消赞操作的错误信息
	err := like.FavouriteAction(userid, videoid, int32(actiontype))
	if err == nil {
		log.Printf("方法like.FavouriteAction(userid, videoid, int32(actiontype) 成功")
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action success",
		})
	} else {
		log.Printf("方法like.FavouriteAction(userid, videoid, int32(actiontype) 失败：%v", err)
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 1,
			StatusMsg:  "favourite action fail",
		})
	}
}

//获取点赞列表;
func GetFavouriteList(c *gin.Context) {
	user_id := c.Query("user_id")
	userid, _ := strconv.ParseInt(user_id, 10, 64)
	like := GetVideo()
	videos, err := like.GetFavouriteList(userid)
	if err == nil {
		log.Printf("方法like.GetFavouriteList(userid) 成功")
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 0,
			StatusMsg:  "get favouritelist success",
			Videolist:  videos,
		})
	} else {
		log.Printf("方法like.GetFavouriteList(userid) 失败：%v", err)
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 1,
			StatusMsg:  "get favouritelist fail ",
		})
	}
}
