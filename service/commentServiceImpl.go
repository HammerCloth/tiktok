package service

import (
	"TikTok/config"
	"TikTok/dao"
	"log"
	"time"
)

type CommentServiceImpl struct {
	UserService
}

//1、使用video id 查询Comment数量
func (c CommentServiceImpl) CountFromVideoId(id int64) (int64, error) {
	return dao.NewCommentDaoInstance().Count(id)
}

//2、发表评论
func (c CommentServiceImpl) Send(comment dao.Comment) error {
	//数据准备
	var commentInfo dao.Comment
	commentInfo.Video_id = comment.Video_id         //评论视频id传入
	commentInfo.User_id = comment.User_id           //评论用户id传入
	commentInfo.Comment_text = comment.Comment_text //评论内容传入
	commentInfo.Cancel = config.ValidComment        //评论状态，0，有效
	nowTime := time.Now().Format(config.DateTime)
	commentInfo.Create_date = nowTime //评论时间记录

	return dao.NewCommentDaoInstance().InsertComment(commentInfo)
}

//3、删除评论，传入评论id
func (c CommentServiceImpl) DelComment(id int64) error {
	return dao.NewCommentDaoInstance().DeleteComment(id)
}

//4、查看评论列表-返回评论list
func (c CommentServiceImpl) GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	//1.先查询评论列表信息
	commentList, err := dao.NewCommentDaoInstance().GetCommentList(videoId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	//提前定义好切片长度
	commentInfoList := make([]CommentInfo, 0, len(commentList))
	for _, comment := range commentList {
		var commentInfo CommentInfo
		commentInfo.Id = comment.Id
		commentInfo.Content = comment.Comment_text
		commentInfo.Create_date = comment.Create_date
		//2.根据查询到的评论用户id和当前用户id，查询评论用户信息
		impl := UserServiceImpl{
			FollowService: &FollowServiceImp{},
		}
		commentInfo.UserInfo, err = impl.GetUserByIdWithCurId(comment.User_id, userId)
		if err != nil {
			return nil, err
		}
		//3.组装list
		commentInfoList = append(commentInfoList, commentInfo)
	}
	return commentInfoList, nil
}
