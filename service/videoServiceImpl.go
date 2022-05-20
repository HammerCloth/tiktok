package service

import (
	"TikTok/config"
	"TikTok/dao"
	"bytes"
	"fmt"
	"github.com/jinzhu/copier"
	"github.com/satori/go.uuid"
	"log"
	"mime/multipart"
	"os/exec"
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
		log.Printf("方法dao.GetVideosByLastTime(lastTime) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法dao.GetVideosByLastTime(lastTime) 成功：%v", tableVideos)
	//将数据通过copyVideos进行处理
	err = videoService.copyVideos(&videos, &tableVideos, userId)
	if err != nil {
		log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 失败：%v", err)
		return nil, time.Time{}, err
	}
	log.Printf("方法videoService.copyVideos(&videos, &tableVideos, userId) 成功")
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
		log.Printf("方法dao.GetVideoByVideoId(videoId) 失败：%v", err)
		return video, err
	}
	log.Printf("方法dao.GetVideoByVideoId(videoId) 成功")

	//将同名字段进行拷贝
	err = copier.Copy(&video, &data)
	if err != nil {
		log.Printf("方法copier.Copy(&video, &data) 失败：%v", err)
		return Video{}, err
	}
	log.Printf("方法copier.Copy(&video, &data) 成功")

	//插入Author
	video.Author, err = videoService.GetUserByIdWithCurId(data.AuthorId, userId)
	if err != nil {
		log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 失败：%v", err)
		return video, err
	}
	log.Printf("方法videoService.GetUserByIdWithCurId(data.AuthorId, userId) 成功")

	//插入点赞数量
	likeCount, err := videoService.FavouriteCount(data.ID)
	if err != nil {
		log.Printf("方法videoService.FavouriteCount(data.ID) 失败：%v", err)
		return video, err
	}
	log.Printf("方法videoService.FavouriteCount(data.ID) 成功")

	video.FavoriteCount = likeCount
	//获取该视屏的评论数字
	commentCount, err := videoService.CountFromVideoId(data.ID)
	if err != nil {
		log.Printf("方法videoService.CountFromVideoId(data.ID) 失败：%v", err)
		return video, err
	}
	log.Printf("方法videoService.CountFromVideoId(data.ID) 成功")

	video.CommentCount = commentCount
	//获取当前用户是否点赞了该视频
	isFavourit, err := videoService.IsFavourit(video.Id, userId)
	if err != nil {
		log.Printf("方法videoService.IsFavourit(video.Id, userId) 失败：%v", err)
	} else {
		log.Printf("方法videoService.IsFavourit(video.Id, userId) 成功")
	}
	video.IsFavorite = isFavourit
	return video, nil
}

// Publish
// 将传入的视频流保存在文件服务器中，并存储在mysql表中
func (videoService *VideoServiceImpl) Publish(data *multipart.FileHeader, userId int64) error {
	//将视频流上传到视频服务器，保存视频链接
	file, err := data.Open()
	if err != nil {
		log.Printf("方法data.Open() 失败%v", err)
		return err
	}
	log.Printf("方法data.Open() 成功")
	//生成一个uuid作为视频的名字
	videoName := uuid.NewV4().String()
	log.Printf("生成视频名称%v", videoName)
	err = dao.VideoFTP(file, videoName)
	if err != nil {
		log.Printf("方法dao.VideoFTP(file, videoName) 失败%v", err)
		return err
	}
	log.Printf("方法dao.VideoFTP(file, videoName) 成功")
	defer file.Close()
	//在服务器上执行ffmpeg 从视频流中获取第一帧截图，并上传图片服务器，保存图片链接
	imageName := uuid.NewV4().String() + ".jpg"
	cmdArguments := []string{"-ss", "00:00:01", "-i", "/home/ftpuser/video/" + videoName + ".mp4", "-vframes", "1", "/home/ftpuser/images/" + imageName}
	cmd := exec.Command("ffmpeg", cmdArguments...)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		log.Printf("ffmpeg 出错！！！！")
		log.Fatal(err)
	}
	fmt.Printf("command output: %q", out.String())
	//组装并持久化
	err = dao.Save(videoName, imageName, userId)
	if err != nil {
		log.Printf("方法dao.Save(videoName, imageName, userId) 失败%v", err)
		return err
	}
	log.Printf("方法dao.Save(videoName, imageName, userId) 成功")
	return nil
}

// List
// 通过userId来查询对应用户发布的视频，并返回对应的视频数组
func (videoService *VideoServiceImpl) List(userId int64) ([]Video, error) {
	//依据用户id查询所有的视频，获取视频列表
	data, err := dao.GetVideosByAuthorId(userId)
	if err != nil {
		log.Printf("方法dao.GetVideosByAuthorId(%v)失败:%v", userId, err)
		return nil, err
	}
	log.Printf("方法dao.GetVideosByAuthorId(%v)成功:%v", userId, data)
	//提前定义好切片长度
	result := make([]Video, 0, len(data))
	//调用拷贝方法，将数据进行转换
	err = videoService.copyVideos(&result, &data, userId)
	if err != nil {
		log.Printf("方法videoService.copyVideos(&result, &data, %v)失败:%v", userId, err)
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
		err := copier.Copy(&video, &temp)
		if err != nil {
			log.Printf("copier.Copy(&video, &temp) 失败：%v", err)
			return err
		}
		log.Println("copier.Copy(&video, &temp) 成功")
		//获取对应的user
		author, err := videoService.GetUserByIdWithCurId(temp.AuthorId, userId)
		if err != nil {
			log.Printf("videoService.GetUserByIdWithCurId(temp.AuthorId, userId) 失败：%v", err)
			return err
		}
		log.Println("videoService.GetUserByIdWithCurId(temp.AuthorId, userId) 成功")
		video.Author = author
		//获取该视屏的点赞数字
		likeCount, err := videoService.FavouriteCount(temp.ID)
		if err != nil {
			log.Printf("videoService.FavouriteCount(temp.ID) 失败：%v", err)
			return err
		}
		log.Printf("videoService.FavouriteCount(temp.ID) 成功")
		video.FavoriteCount = likeCount
		//获取该视屏的评论数字
		commentCount, err := videoService.CountFromVideoId(temp.ID)
		if err != nil {
			log.Printf("videoService.CountFromVideoId(temp.ID) 失败：%v", err)
			return err
		}
		log.Printf("videoService.CountFromVideoId(temp.ID) 成功")
		video.CommentCount = commentCount
		//获取当前用户是否点赞了该视频
		isFavourit, err := videoService.IsFavourit(video.Id, userId)
		if err != nil {
			log.Printf("videoService.IsFavourit(video.Id, userId) 失败：%v", err)
		} else {
			log.Printf("videoService.IsFavourit(video.Id, userId) 成功")
		}
		video.IsFavorite = isFavourit
		*result = append(*result, video)

	}
	return nil
}
