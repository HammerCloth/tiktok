package dao

import (
	"TikTok/config"
	"errors"
	"github.com/jinzhu/gorm"
	"log"
	"sync"
)

//评论信息-数据库中的结构体-dao层使用
type Comment struct {
	Id           int64  //评论id
	User_id      int64  //评论用户id
	Video_id     int64  //视频id
	Comment_text string //评论内容
	Create_date  string //评论发布的日期mm-dd
	Cancel       int32  //取消评论为1，发布评论为0
}

// TableName 修改表名映射
func (Comment) TableName() string {
	return "comments"
}

type CommentDao struct {
}

var (
	commentDao  *CommentDao
	commentOnce sync.Once
)

func NewCommentDaoInstance() *CommentDao {
	commentOnce.Do(
		func() {
			commentDao = &CommentDao{}
		})
	return commentDao
}

//1、使用video id 查询Comment数量
func (*CommentDao) Count(video_id int64) (int64, error) {
	Init()
	var count int64
	err := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": video_id, "cancel": config.ValidComment}).Count(&count).Error
	if err != nil {
		return -1, errors.New("find comments count failed")
	}
	return count, nil
}

//2、发表评论
func (*CommentDao) InsertComment(comment Comment) error {
	Init()
	err := Db.Model(Comment{}).Create(&comment).Error
	if err != nil {
		return errors.New("create comment failed")
	}
	return nil
}

//3、删除评论，传入评论id
func (*CommentDao) DeleteComment(id int64) error {
	Init()
	var commentInfo Comment
	//先查询是否有此评论（正常肯定是有的吧）
	result := Db.Model(Comment{}).Where(map[string]interface{}{"id": id, "cancel": config.ValidComment}).First(&commentInfo)
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		return errors.New("del comment is not exist")
	}
	err := Db.Model(Comment{}).Where("id = ?", id).Update("cancel", config.InvalidComment).Error
	if err != nil {
		return errors.New("del comment failed")
	}
	return nil
}

//4.根据视频id查询所属评论全部列表信息
func (*CommentDao) GetCommentList(videoId int64) ([]Comment, error) {
	Init()
	var commentList []Comment
	result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).Find(&commentList)
	if result.RowsAffected == 0 {
		return commentList, errors.New("there are no comments")
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
		return commentList, errors.New("get commentList failed")
	}
	return commentList, nil
}
