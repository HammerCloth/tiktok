package controller

import (
	"TikTok/service"
	"github.com/gin-gonic/gin"
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
	if like.FavouriteAction(userid, videoid, int32(actiontype)) == nil {
		c.JSON(http.StatusOK, likeResponse{
			StatusCode: 0,
			StatusMsg:  "favourite action success",
		})
	} else {
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
	like := new(service.LikeServiceImpl)
	videos, err := like.GetFavouriteList(userid)
	if err == nil {
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 0,
			StatusMsg:  "get favouritelist success",
			Videolist:  videos,
		})
	} else {
		c.JSON(http.StatusOK, GetFavouriteListResponse{
			StatusCode: 1,
			StatusMsg:  "can't find favouritelist ",
		})
	}
}
