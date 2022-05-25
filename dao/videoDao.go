package dao

import (
	"TikTok/config"
	"github.com/dutchcoders/goftp"
	"io"
	"time"
)

type TableVideo struct {
	ID          int64
	AuthorId    int64 `copier:"-"` //在拷贝时忽略
	PlayUrl     string
	CoverUrl    string
	PublishTime time.Time `copier:"-"` //在拷贝时忽略
	Title       string    //视频名，5.23添加
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
	tableVideo.ID = videoId
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
	//Init()
	result := Db.Where("publish_time<?", lastTime).Order("publish_time desc").Limit(config.VideoCount).Find(&videos)
	if result.Error != nil {
		return videos, result.Error
	}
	return videos, nil
}

// VideoFTP
// 通过ftp将视频传入服务器
func VideoFTP(file io.Reader, videoName string) error {
	//初始化ftp
	ftp, err := initFTP()
	if err != nil {
		return err
	}
	//转到video相对路线下
	err = ftp.Cwd("video")
	if err != nil {
		return err
	}
	if err := ftp.Stor(videoName+".mp4", file); err != nil {
		return err
	}
	return nil
}

//初始化FTP
func initFTP() (*goftp.FTP, error) {
	//获取到ftp的链接
	connect, err := goftp.Connect(config.ConConfig)
	if err != nil {
		return nil, err
	}
	//登录
	err = connect.Login(config.FtpUser, config.FtpPsw)
	if err != nil {
		return nil, err
	}
	return connect, nil
}

// ImageFTP
// 将图片传入FTP服务器中，但是这里要注意图片的格式随着名字一起给,同时调用时需要自己结束流
func ImageFTP(file io.Reader, imageName string) error {
	//初始化ftp
	ftp, err := initFTP()
	if err != nil {
		return err
	}
	//转到video相对路线下
	err = ftp.Cwd("images")
	if err != nil {
		return err
	}
	if err := ftp.Stor(imageName, file); err != nil {
		return err
	}
	return nil
}

// Save 保存视频记录
func Save(videoName string, imageName string, authorId int64, title string) error {
	//Init()
	var video TableVideo
	video.PublishTime = time.Now()
	video.PlayUrl = config.PlayUrlPrefix + videoName + ".mp4"
	video.CoverUrl = config.CoverUrlPrefix + imageName
	video.AuthorId = authorId
	video.Title = title
	result := Db.Save(&video)
	if result.Error != nil {
		return result.Error
	}
	return nil
}
