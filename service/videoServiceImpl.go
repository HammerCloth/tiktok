package service

import (
	"TikTok/dao"
	"github.com/jinzhu/copier"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
	UserService
}

// Feed
// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error) {
	return nil, time.Time{}, nil
}

// GetVideo
// 传入视频id获得对应的视频对象
func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	var video Video
	data, err := dao.GetVideoByVideoId(videoId)
	if err != nil {
		return video, err
	}
	copier.Copy(&video, &data)
	video.Author, err = videoService.GetUserByIdWithCurId(data.AuthorId, userId)
	if err != nil {
		return video, err
	}
	return video, nil
}

// Publish
// 将传入的视频流保存在文件服务器中，并存储在mysql表中
func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64) error {
	return nil
}

// List
// 通过userId来查询对应用户发布的视频，并返回对应的视频数组
func (videoService *VideoServiceImpl) List(userId int64) ([]Video, error) {
	//依据用户id查询所有的视频，获取视频列表
	data, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		return nil, err
	}
	//提前定义好切片长度
	result := make([]Video, 0, len(data))
	for _, temp := range data {
		var video Video
		//进行拷贝操作
		copier.Copy(&video, &temp)
		//获取对应的user
		video.Author, err = videoService.GetUserByIdWithCurId(temp.AuthorId, temp.AuthorId)
		if err != nil {
			return nil, err
		}
		result = append(result, video)
	}
	return result, nil
}
