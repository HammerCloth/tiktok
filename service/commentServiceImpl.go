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
	log.Println("CommentService-Send: running") //函数已运行
	//数据准备
	var commentInfo dao.Comment
	commentInfo.Video_id = comment.Video_id         //评论视频id传入
	commentInfo.User_id = comment.User_id           //评论用户id传入
	commentInfo.Comment_text = comment.Comment_text //评论内容传入
	commentInfo.Cancel = config.ValidComment        //评论状态，0，有效
	//nowTime := time.Now().Format(config.DateTime)
	commentInfo.Create_date = time.Now() //评论时间记录

	return dao.NewCommentDaoInstance().InsertComment(commentInfo)
}

//3、删除评论，传入评论id
func (c CommentServiceImpl) DelComment(id int64) error {
	log.Println("CommentService-DelComment: running") //函数已运行
	return dao.NewCommentDaoInstance().DeleteComment(id)
}

//4、查看评论列表-返回评论list
func (c CommentServiceImpl) GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	log.Println("CommentService-GetList: running") //函数已运行
	/*
		//法一、使用SQL语句查询评论列表及用户信息，但是好像必须嵌套user信息，直接查的就不行
		commentInfoList := make([]CommentInfo, 1)
		err := dao.Db.Raw("select T.cid,T.user_id,T.`name`,T.follow_cnt,T.follower_cnt,"+
			"\nif(f.cancel is null,'false','true') isFollow,"+
			"\nT.comment_text,T.create_date"+
			"\nfrom follows f right join\n("+
			"\n\tselect cid,vid,id user_id,`name`,comment_text,create_date,"+
			"\n\tcount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
			"\n\tcount(if(tag = 'follow' and cancel is not null,1,null)) follow_count,"+
			"\n\tfrom"+
			"\n\t\t("+
			"\n\t\tselect c.id cid,u.id,c.video_id vid,`name`,f.cancel,comment_text,create_date,'follower' tag"+
			"\n\t\tfrom comments c join users u on c.user_id = u.id and c.cancel = 0"+
			"\n\t\tleft join follows f on u.id = f.user_id and f.cancel = 0"+
			"\n\t\tunion all"+
			"\n\t\tselect c.id cid,u.id,c.video_id vid,`name`,f.cancel,comment_text,create_date,'follow' tag"+
			"\n\t\tfrom comments c join users u on c.user_id = u.id and c.cancel = 0"+
			"\n\t\tleft join follows f on u.id = f.follower_id and f.cancel = 0"+
			"\n\t\t) T"+
			"\n\tgroup by cid,vid,id,`name`,comment_text,create_date"+
			"\n) T on f.follower_id = T.user_id and f.cancel = 0 and f.user_id = ?"+
			"\nwhere vid = ? group by cid order by create_date desc", userId, videoId).Scan(&commentInfoList).Error

		if nil != err {
			log.Println("CommentService-GetList: sql error") //sql查询出错
			return nil, err
		}
		log.Println("CommentService-GetList: get list success") //成功查询到评论列表
		return commentInfoList, nil
	*/

	//法二：调用dao，先查评论，再循环查用户信息：
	//1.先查询评论列表信息
	commentList, err := dao.NewCommentDaoInstance().GetCommentList(videoId)
	if err != nil {
		log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
		return nil, err
	}
	//提前定义好切片长度
	commentInfoList := make([]CommentInfo, 0, len(commentList))
	for _, comment := range commentList {
		var commentInfo CommentInfo
		commentInfo.Id = comment.Id
		commentInfo.Content = comment.Comment_text
		//时间解析
		//timeStr:=comment.Create_date
		commentInfo.Create_date = comment.Create_date.Format(config.DateTime)
		//dateTime, err := time.Parse(config.DateTime, comment.Create_date)
		//commentInfo.Create_date = dateTime.Format(config.DateTime)
		log.Printf("commentTime1:%v\n", comment.Create_date)
		log.Printf("commentTime2:%v", commentInfo.Create_date)
		//2.根据查询到的评论用户id和当前用户id，查询评论用户信息
		impl := UserServiceImpl{
			FollowService: &FollowServiceImp{},
		}
		commentInfo.UserInfo, err = impl.GetUserByIdWithCurId(comment.User_id, userId)
		if err != nil {
			log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
			return nil, err
		}
		//3.组装list
		commentInfoList = append(commentInfoList, commentInfo)
	}
	log.Println("CommentService-GetList: return list success") //函数执行成功，返回正确信息
	return commentInfoList, nil
}
