package service

import (
	"mime/multipart"
	"time"
)

type Video struct {
	Id            int64  `copier:"ID" json:"id,omitempty"` //指定别名
	Author        User   `copier:"-" json:"author"`        //在拷贝时忽略
	PlayUrl       string `json:"play_url,omitempty"`
	CoverUrl      string `json:"cover_url,omitempty"`
	FavoriteCount int64  `json:"favorite_count,omitempty"`
	CommentCount  int64  `json:"comment_count,omitempty"`
	IsFavorite    bool   `json:"is_favorite,omitempty"`
}

type VideoService interface {
	// Feed
	// 通过传入时间戳，当前用户的id，返回对应的视频切片数组，以及视频数组中最早的发布时间
	Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error)

	// GetVideo
	// 传入视频id获得对应的视频对象
	GetVideo(videoId int64, userId int64) (Video, error)

	// Publish
	// 将传入的视频流保存在文件服务器中，并存储在mysql表中
	Publish(data *multipart.FileHeader, userId int64) error

	// List
	// 通过userId来查询对应用户发布的视频，并返回对应的视频切片数组
	List(userId int64) ([]Video, error)
}
