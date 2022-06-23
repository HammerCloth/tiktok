package dao

import (
	"TikTok/config"
	"errors"
	"log"
	"time"
)

// Comment
// 评论信息-数据库中的结构体-dao层使用
type Comment struct {
	Id          int64     //评论id
	UserId      int64     //评论用户id
	VideoId     int64     //视频id
	CommentText string    //评论内容
	CreateDate  time.Time //评论发布的日期mm-dd
	Cancel      int32     //取消评论为1，发布评论为0
}

// TableName 修改表名映射
func (Comment) TableName() string {
	return "comments"
}

// Count
// 1、使用video id 查询Comment数量
func Count(videoId int64) (int64, error) {
	log.Println("CommentDao-Count: running") //函数已运行
	//Init()
	var count int64
	//数据库中查询评论数量
	err := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).Count(&count).Error
	if err != nil {
		log.Println("CommentDao-Count: return count failed") //函数返回提示错误信息
		return -1, errors.New("find comments count failed")
	}
	log.Println("CommentDao-Count: return count success") //函数执行成功，返回正确信息
	return count, nil
}

//CommentIdList 根据视频id获取评论id 列表
func CommentIdList(videoId int64) ([]string, error) {
	var commentIdList []string
	err := Db.Model(Comment{}).Select("id").Where("video_id = ?", videoId).Find(&commentIdList).Error
	if err != nil {
		log.Println("CommentIdList:", err)
		return nil, err
	}
	return commentIdList, nil
}

// InsertComment
// 2、发表评论
func InsertComment(comment Comment) (Comment, error) {
	log.Println("CommentDao-InsertComment: running") //函数已运行
	//数据库中插入一条评论信息
	err := Db.Model(Comment{}).Create(&comment).Error
	if err != nil {
		log.Println("CommentDao-InsertComment: return create comment failed") //函数返回提示错误信息
		return Comment{}, errors.New("create comment failed")
	}
	log.Println("CommentDao-InsertComment: return success") //函数执行成功，返回正确信息
	return comment, nil
}

// DeleteComment
// 3、删除评论，传入评论id
func DeleteComment(id int64) error {
	log.Println("CommentDao-DeleteComment: running") //函数已运行
	var commentInfo Comment
	//先查询是否有此评论
	result := Db.Model(Comment{}).Where(map[string]interface{}{"id": id, "cancel": config.ValidComment}).First(&commentInfo)
	if result.RowsAffected == 0 { //查询到此评论数量为0则返回无此评论
		log.Println("CommentDao-DeleteComment: return del comment is not exist") //函数返回提示错误信息
		return errors.New("del comment is not exist")
	}
	//数据库中删除评论-更新评论状态为-1
	err := Db.Model(Comment{}).Where("id = ?", id).Update("cancel", config.InvalidComment).Error
	if err != nil {
		log.Println("CommentDao-DeleteComment: return del comment failed") //函数返回提示错误信息
		return errors.New("del comment failed")
	}
	log.Println("CommentDao-DeleteComment: return success") //函数执行成功，返回正确信息
	return nil
}

// GetCommentList
// 4.根据视频id查询所属评论全部列表信息
func GetCommentList(videoId int64) ([]Comment, error) {
	log.Println("CommentDao-GetCommentList: running") //函数已运行
	//数据库中查询评论信息list
	var commentList []Comment
	result := Db.Model(Comment{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.ValidComment}).
		Order("create_date desc").Find(&commentList)
	//若此视频没有评论信息，返回空列表，不报错
	if result.RowsAffected == 0 {
		log.Println("CommentDao-GetCommentList: return there are no comments") //函数返回提示无评论
		return nil, nil
	}
	//若获取评论列表出错
	if result.Error != nil {
		log.Println(result.Error.Error())
		log.Println("CommentDao-GetCommentList: return get comment list failed") //函数返回提示获取评论错误
		return commentList, errors.New("get comment list failed")
	}
	log.Println("CommentDao-GetCommentList: return commentList success") //函数执行成功，返回正确信息
	return commentList, nil
}
