package dao

import (
	"TikTok/config"
	"TikTok/middleware/ftp"
	"io"
	"log"
	"time"
)

type TableVideo struct {
	Id          int64 `json:"id"`
	AuthorId    int64
	PlayUrl     string `json:"play_url"`
	CoverUrl    string `json:"cover_url"`
	PublishTime time.Time
	Title       string `json:"title"` //视频名，5.23添加
}

// TableName
//	将TableVideo映射到videos，
//	这样我结构体到名字就不需要是Video了，防止和我Service层到结构体名字冲突
func (TableVideo) TableName() string {
	return "videos"
}

// GetVideosByAuthorId
// 根据作者的id来查询对应数据库数据，并TableVideo返回切片
func GetVideosByAuthorId(authorId int64) ([]TableVideo, error) {
	//建立结果集接收
	var data []TableVideo
	//初始化db
	//Init()
	result := Db.Where(&TableVideo{AuthorId: authorId}).Find(&data)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return data, nil
}

// GetVideoByVideoId
// 依据VideoId来获得视频信息
func GetVideoByVideoId(videoId int64) (TableVideo, error) {
	var tableVideo TableVideo
	tableVideo.Id = videoId
	//Init()
	result := Db.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}

// GetVideosByLastTime
// 依据一个时间，来获取这个时间之前的一些视频
func GetVideosByLastTime(lastTime time.Time) ([]TableVideo, error) {
	videos := make([]TableVideo, config.VideoCount)
	result := Db.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

// VideoFTP
// 通过ftp将视频传入服务器
func VideoFTP(file io.Reader, videoName string) error {
	//转到video相对路线下
	err := ftp.MyFTP.Cwd("video")
	if err != nil {
		log.Println("转到路径video失败！！！")
	} else {
		log.Println("转到路径video成功！！！")
	}
	err = ftp.MyFTP.Stor(videoName+".mp4", file)
	if err != nil {
		log.Println("上传视频失败！！！！！")
		return err
	}
	log.Println("上传视频成功！！！！！")
	return nil
}

// ImageFTP
// 将图片传入FTP服务器中，但是这里要注意图片的格式随着名字一起给,同时调用时需要自己结束流
func ImageFTP(file io.Reader, imageName string) error {
	//转到video相对路线下
	err := ftp.MyFTP.Cwd("images")
	if err != nil {
		log.Println("转到路径images失败！！！")
		return err
	}
	log.Println("转到路径images成功！！！")
	if err = ftp.MyFTP.Stor(imageName, file); err != nil {
		log.Println("上传图片失败！！！！！")
		return err
	}
	log.Println("上传图片成功！！！！！")
	return nil
}

// Save 保存视频记录
func Save(videoName string, imageName string, authorId int64, title string) error {
	var video TableVideo
	video.PublishTime = time.Now()
	video.PlayUrl = config.PlayUrlPrefix + videoName + ".mp4"
	video.CoverUrl = config.CoverUrlPrefix + imageName + ".jpg"
	video.AuthorId = authorId
	video.Title = title
	result := Db.Save(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

// GetVideoIdsByAuthorId
// 通过作者id来查询发布的视频id切片集合
func GetVideoIdsByAuthorId(authorId int64) ([]int64, error) {
	var id []int64
	//通过pluck来获得单独的切片
	result := Db.Model(&TableVideo{}).Where("author_id", authorId).Pluck("id", &id)
	//如果出现问题，返回对应到空，并且返回error
	if result.Error != nil {
		return nil, result.Error
	}
	return id, nil
}
