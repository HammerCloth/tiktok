package service

import (
	"TikTok/config"
	"TikTok/dao"
	"github.com/jinzhu/copier"
	"mime/multipart"
	"time"
)

type VideoServiceImpl struct {
	UserService
	LikeService
	CommentService
}

// Feed
// 通过传入时间戳，当前用户的id，返回对应的视频数组，以及视频数组中最早的发布时间
func (videoService VideoServiceImpl) Feed(lastTime time.Time, userId int64) ([]Video, time.Time, error) {
	//创建对应返回视频的切片数组
	videos := make([]Video, 0, config.VideoCount)
	//根据传入的时间，获得传入时间前n个视频，可以通过config.videoCount来控制
	tableVideos, err := dao.GetVideosByLastTime(lastTime)
	if err != nil {
		return nil, time.Time{}, err
	}
	//将数据通过copyVideos进行处理
	err = videoService.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		return nil, time.Time{}, err
	}
	//返回数据，同时获得视频中最早的时间返回
	return videos, tableVideos[config.VideoCount-1].PublishTime, nil
}

// GetVideo
// 传入视频id获得对应的视频对象，注意还需要传入当前的用户id
func (videoService *VideoServiceImpl) GetVideo(videoId int64, userId int64) (Video, error) {
	//初始化video对象
	var video Video
	//从数据库中查询数据
	data, err := dao.GetVideoByVideoId(videoId)
	if err != nil {
		return video, err
	}
	//将同名字段进行拷贝
	copier.Copy(&video, &data)
	//插入Author
	video.Author, err = videoService.GetUserByIdWithCurId(data.AuthorId, userId)
	if err != nil {
		return video, err
	}
	//插入点赞数量
	likeCount, err := videoService.FavouriteCount(data.ID)
	if err != nil {
		return video, err
	}
	video.FavoriteCount = likeCount
	//获取该视屏的评论数字
	commentCount, err := videoService.CountFromVideoId(data.ID)
	if err != nil {
		return video, err
	}
	video.CommentCount = commentCount
	//获取当前用户是否点赞了该视频
	isFavourit, err := videoService.IsFavourit(video.Id, userId)
	if err != nil {
		return video, err
	}
	video.IsFavorite = isFavourit
	return video, nil
}

// Publish
// 将传入的视频流保存在文件服务器中，并存储在mysql表中
func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64) error {
	//todo 从视频流中获取第一帧截图，并上传图片服务器，保存图片链接
	//todo 将视频流上传到视频服务器，保存视频链接
	//todo 组装并持久化
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
	//调用拷贝方法，将数据进行转换
	err = videoService.copyVideos(&result, &data, userId)
	if err != nil {
		return nil, err
	}
	//如果数据没有问题，则直接返回
	return result, nil
}

// 该方法可以将数据进行拷贝和转换，并从其他方法获取对应的数据
func (videoService *VideoServiceImpl) copyVideos(result *[]Video, data *[]dao.TableVideo, userId int64) error {
	for _, temp := range *data {
		var video Video
		//进行拷贝操作
		copier.Copy(&video, &temp)
		//获取对应的user
		author, err := videoService.GetUserByIdWithCurId(temp.AuthorId, userId)
		if err != nil {
			return err
		}
		video.Author = author
		//获取该视屏的点赞数字
		likeCount, err := videoService.FavouriteCount(temp.ID)
		if err != nil {
			return err
		}
		video.FavoriteCount = likeCount
		//获取该视屏的评论数字
		commentCount, err := videoService.CountFromVideoId(temp.ID)
		if err != nil {
			return err
		}
		video.CommentCount = commentCount
		//获取当前用户是否点赞了该视频
		isFavourit, err := videoService.IsFavourit(video.Id, userId)
		if err != nil {
			return err
		}
		video.IsFavorite = isFavourit
		*result = append(*result, video)

	}
	return nil
}
