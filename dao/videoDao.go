package dao

import (
	"time"
)

type TableVideo struct {
	ID          int64
	AuthorId    int64 `copier:"-"` //在拷贝时忽略
	PlayUrl     string
	CoverUrl    string
	PublishTime time.Time `copier:"-"` //在拷贝时忽略
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
	Init()
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
	Init()
	result := Db.First(&tableVideo)
	if result.Error != nil {
		return tableVideo, result.Error
	}
	return tableVideo, nil

}
