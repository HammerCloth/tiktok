package dao

import (
	"TikTok/config"
	"errors"
	"log"
	"sync"
)

//like表的结构。
type Like struct {
	Id       int64 //自增主键
	User_id  int64 //点赞用户id
	Video_id int64 //视频id
	Cancel   int8  //是否点赞，0为点赞，1为取消赞
}

// TableName 修改表名映射
func (Like) TableName() string {
	return "likes"
}

type LikeDao struct {
}

var (
	likeDao  *LikeDao
	likeOnce sync.Once
)

func NewLikeDaoInstance() *LikeDao {
	likeOnce.Do(
		func() {
			likeDao = &LikeDao{}
		})
	return likeDao
}

//1.根据videoid获取点赞数量
func (*LikeDao) GetLikeCount(videoId int64) (int64, error) {
	Init()
	var count int64
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.Islike}).
		Count(&count).Error
	if err != nil {
		return 0, errors.New("An unknown exception occurred in the query")
	} else {
		return count, nil
	}
}

//2.根据userid，videoid,action_type点赞或者取消赞
func (*LikeDao) UpdateLike(userId int64, videoId int64, action_type int32) error {
	Init()
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", action_type).Error
	if err != nil {
		return errors.New("update data fail")
	}
	return nil
}

//3、插入点赞数据
func (*LikeDao) InsertLike(likedata Like) error {
	Init()
	err := Db.Model(Like{}).Create(&likedata).Error
	if err != nil {
		return errors.New("insert data fail")
	}
	return nil
}

//4.根据userid,videoid查询点赞信息
func (*LikeDao) GetLikeInfo(userId int64, videoId int64) (Like, error) {
	Init()
	var likeInfo Like
	result := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&likeInfo)
	//如果获取likeInfo失败
	if result.Error != nil {
		log.Println(result.Error.Error())
		return likeInfo, errors.New("get likeInfo failed")
	}
	//查询数据为0
	if result.RowsAffected == 0 {
		log.Println(result.Error.Error())
		return likeInfo, errors.New("can't find this data")
	}
	return likeInfo, nil
}

//5.根据userid查询所属点赞全部列表信息
func (*LikeDao) GetLikeList(userId int64) ([]Like, error) {
	Init()
	var likeList []Like
	result := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "cancel": config.Islike}).
		Find(&likeList)
	//如果获取likelist失败
	if result.Error != nil {
		log.Println(result.Error.Error())
		return likeList, errors.New("get likeList failed")
	}
	//查询数据为0
	if result.RowsAffected == 0 {
		return likeList, errors.New("there are no likes")
	}
	return likeList, nil
}
