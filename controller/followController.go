package controller

import (
	"TikTok/service"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// RelationActionResp 关注和取消关注需要返回结构。
type RelationActionResp struct {
	Response
}

// FollowingResp 获取关注列表需要返回的结构。
type FollowingResp struct {
	Response
	users []service.User `json:"users"`
}

// FollowersResp 获取粉丝列表需要返回的结构。
type FollowersResp struct {
	Response
	users []service.User `json:"users"`
}

// RelationAction 处理关注和取消关注请求。
func RelationAction(c *gin.Context) {
	userId, err1 := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	toUserId, err2 := strconv.ParseInt(c.GetString("to_user_id"), 10, 64)
	actionType, err3 := strconv.ParseInt(c.GetString("action_type"), 10, 64)
	// 传入参数格式有问题。
	if nil != err1 || nil != err2 || nil != err3 || actionType < 1 || actionType > 2 {
		c.JSON(http.StatusOK, RelationActionResp{
			Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误",
			},
		})
		return
	}
	// 正常处理
	fsi := service.NewFSIInstance()
	switch {
	// 关注
	case 1 == actionType:
		fsi.AddFollowRelation(userId, toUserId)
	// 取关
	case 2 == actionType:
		fsi.DeleteFollowRelation(userId, toUserId)
	}
	c.JSON(http.StatusOK, RelationActionResp{
		Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
	})
}

// GetFollowing 处理获取关注列表请求
func GetFollowing(c *gin.Context) {
	userId, err := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			users: nil,
		})
		return
	}
	// 正常获取关注列表
	fsi := service.NewFSIInstance()
	users, err := fsi.GetFollowing(userId)
	// 获取关注列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取关注列表时出错。",
			},
			users: nil,
		})
		return
	}
	// 成功获取到关注列表。
	c.JSON(http.StatusOK, FollowingResp{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		users: users,
	})
}

// GetFollowers 处理获取关注列表请求
func GetFollowers(c *gin.Context) {
	userId, err := strconv.ParseInt(c.GetString("user_id"), 10, 64)
	// 用户id解析出错。
	if nil != err {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "用户id格式错误。",
			},
			users: nil,
		})
		return
	}
	// 正常获取粉丝列表
	fsi := service.NewFSIInstance()
	users, err := fsi.GetFollowers(userId)
	// 获取关注列表时出错。
	if err != nil {
		c.JSON(http.StatusOK, FollowingResp{
			Response: Response{
				StatusCode: -1,
				StatusMsg:  "获取粉丝列表时出错。",
			},
			users: nil,
		})
		return
	}
	// 成功获取到粉丝列表。
	c.JSON(http.StatusOK, FollowingResp{
		Response: Response{
			StatusCode: 0,
			StatusMsg:  "OK",
		},
		users: users,
	})
}
