package dao

import (
	"TikTok/config"
	"errors"
	"log"
	"sync"
	"time"
)

//评论信息-数据库中的结构体-dao层使用
type Comment struct {
	Id           int64     //评论id
	User_id      int64     //评论用户id
	Video_id     int64     //视频id
	Comment_text string    //评论内容
	Create_date  time.Time //评论发布的日期mm-dd
	Cancel       int32     //取消评论为1，发布评论为0
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
	log.Println("CommentDao-Count: running") //函数已运行
	//Init()
	var count int64
	err := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": video_id, "cancel": config.ValidComment}).Count(&count).Error
	if err != nil {
		log.Println("CommentDao-Count: return count failed") //函数返回提示错误信息
		return -1, errors.New("find comments count failed")
	}
	log.Println("CommentDao-Count: return count success") //函数执行成功，返回正确信息
	return count, nil
}

//2、发表评论
func (*CommentDao) InsertComment(comment Comment) error {
	log.Println("CommentDao-InsertComment: running") //函数已运行
	//Init()
	err := Db.Model(Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") //函数返回提示错误信息
		return errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success") //函数执行成功，返回正确信息
	return nil
}

//3、删除评论，传入评论id
func (*CommentDao) DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running") //函数已运行
	//Init()
	var commentInfo Comment
	//先查询是否有此评论（正常肯定是有的吧）
	result := Db.Model(Comment{}).Where(map[string]interface{}{"id": id, "cancel": config.ValidComment}).First(&commentInfo)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-DeleteComment: return del comment is not exist") //函数返回提示错误信息
		return errors.New("del comment is not exist")
	}
	err := Db.Model(Comment{}).Where("id = ?", id).Update("cancel", config.InvalidComment).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") //函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") //函数执行成功，返回正确信息
	return nil
}

//4.根据视频id查询所属评论全部列表信息
func (*CommentDao) GetCommentList(videoId int64) ([]Comment, error) {
	log.Println("CommentDao-GetCommentList: running") //函数已运行
	//Init()
	var commentList []Comment
	result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).Find(&commentList)
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments") //函数返回提示无评论
		return commentList, errors.New("there are no comments")
	}
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed") //函数返回提示无评论
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success") //函数执行成功，返回正确信息
	return commentList, nil
}
