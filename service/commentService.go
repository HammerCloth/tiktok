package service

import (
	"TikTok/dao"
	"time"
)

// CommentService 接口定义
// 发表评论-使用的结构体-service层引用dao层↑的Comment。
type CommentService interface {
	/*
		一、其他同学（video）需要使用的方法：
	*/
	// CountFromVideoId
	// 1.根据videoId获取视频评论数量的接口
	CountFromVideoId(id int64) (int64, error)

	/*
		二、评论模块自己request实现的方法：
	*/
	// Send
	// 2、发表评论，传进来评论的基本信息，返回保存是否成功的状态描述
	Send(comment dao.Comment) (CommentInfo, error)
	// DelComment
	// 3、删除评论，传入评论id即可，返回错误状态信息
	DelComment(commentId int64) error
	// GetList
	// 4、查看评论列表-返回评论list-在controller层再封装外层的状态信息
	GetList(videoId int64, userId int64) ([]CommentInfo, error)
}

// CommentInfo 查看评论-传出的结构体-service
type CommentInfo struct {
	Id         int64  `json:"id,omitempty"`
	UserInfo   User   `json:"user,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDate string `json:"create_date,omitempty"`
}

type CommentData struct {
	Id            int64     `json:"id,omitempty"`
	UserId        int64     `json:"user_id,omitempty"`
	Name          string    `json:"name,omitempty"`
	FollowCount   int64     `json:"follow_count"`
	FollowerCount int64     `json:"follower_count"`
	IsFollow      bool      `json:"is_follow"`
	Content       string    `json:"content,omitempty"`
	CreateDate    time.Time `json:"create_date,omitempty"`
}
