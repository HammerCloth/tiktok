package service

import (
	"TikTok/dao"
	"errors"
)

type CommentServiceImpl struct {
}

//1、使用video id 查询Comment数量
func (c CommentServiceImpl) CountFromVideoId(id int64) (int64, error) {
	var count int64
	if err := dao.Db.Table("comments").Where("id = ?", id).First(&count); err != nil {
		return -1, errors.New("can't find commentsCount")
	}
	return count, nil
}

//2、发表评论
func (c CommentServiceImpl) Send(comment *Comment) error {
	if comment.Id == 0 {
		return errors.New("Comment id = 0")
	}
	return nil
}

//3、删除评论，传入评论id
func (c CommentServiceImpl) DelComment(id int64) error {
	if id == 0 {
		return errors.New("Comment id = 0")
	}
	return nil
}

//4、查看评论列表-返回评论list
func (c CommentServiceImpl) GetList(vedioId int64, userId int64) ([]CommentInfo, error) {
	//comment := new(dao.Comment)
	if vedioId == 0 || userId == 0 {
		return nil, errors.New("id = 0")
	}
	return nil, nil
}
