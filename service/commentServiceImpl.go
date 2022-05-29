package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware"
	"log"
	"strconv"
)

type CommentServiceImpl struct {
	UserService
}

// CountFromVideoId
// 1、使用video id 查询Comment数量
func (c CommentServiceImpl) CountFromVideoId(videoId int64) (int64, error) {
	//先在缓存中查
	cnt, err := middleware.RdbVCid.ZCard(middleware.Ctx, strconv.FormatInt(videoId, 10)).Result()
	if err != nil {
		return 0, err
	}
	if cnt != 0 {
		return cnt, nil
	}
	//缓存中查不到则去数据库查
	cntDao, err1 := dao.Count(videoId)
	if err1 != nil {
		return 0, nil
	}
	//将评论id切片存入redis，
	go func() {
		//查询评论id list
		cList, _ := dao.CommentIdList(videoId)
		//评论id循环存入redis，
		for _, _comment := range cList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), _comment)
		}
		log.Println("count comment save ids in redis")
	}()
	//返回结果
	return cntDao, nil
}

// Send
// 2、发表评论
func (c CommentServiceImpl) Send(comment dao.Comment) (CommentInfo, error) {
	log.Println("CommentService-Send: running") //函数已运行
	//数据准备
	var commentInfo dao.Comment
	commentInfo.VideoId = comment.VideoId         //评论视频id传入
	commentInfo.UserId = comment.UserId           //评论用户id传入
	commentInfo.CommentText = comment.CommentText //评论内容传入
	commentInfo.Cancel = config.ValidComment      //评论状态，0，有效
	commentInfo.CreateDate = comment.CreateDate   //评论时间记录

	//1.评论信息存储：
	commentRtn, err := dao.InsertComment(commentInfo)
	if err != nil {
		return CommentInfo{}, err
	}
	//2.查询用户信息
	impl := UserServiceImpl{
		FollowService: &FollowServiceImp{},
	}
	userData, err2 := impl.GetUserByIdWithCurId(comment.UserId, comment.UserId)
	if err2 != nil {
		return CommentInfo{}, err2
	}
	//3.拼接
	commentData := CommentInfo{
		Id:         commentRtn.Id,
		UserInfo:   userData,
		Content:    commentRtn.CommentText,
		CreateDate: commentRtn.CreateDate.Format(config.DateTime),
	}
	//将此发表的评论id存入redis-不用mq
	//middleware.MqCommentAdd.CommentPublish(msg.String())
	go func() {
		insertRedisVideoCommentId(strconv.Itoa(int(comment.VideoId)), strconv.Itoa(int(commentRtn.Id)))
		log.Println("send comment save in redis")
	}()
	//返回结果
	return commentData, nil
}

// DelComment
// 3、删除评论，传入评论id
func (c CommentServiceImpl) DelComment(commentId int64) error {
	log.Println("CommentService-DelComment: running") //函数已运行
	//1.先查询redis，若有则删除，返回客户端-再go协程删除数据库-不用mq、考虑没有大量删除的情况.
	//无则在数据库中删除，返回客户端
	n, err := middleware.RdbCVid.Exists(middleware.Ctx, strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		log.Println(err)
	}
	if n > 0 { //在内存中，则找出来删除，然后返回
		vid, err1 := middleware.RdbCVid.Get(middleware.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err1 != nil { //没找到，返回err
			log.Println("redis find CV err:", err1)
		}
		//删除，两个redis都要删除
		del1, err2 := middleware.RdbCVid.Del(middleware.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err2 != nil {
			log.Println(err2)
		}
		del2, err3 := middleware.RdbVCid.SRem(middleware.Ctx, vid, strconv.FormatInt(commentId, 10)).Result()
		if err3 != nil {
			log.Println(err3)
		}
		log.Println("del comment in Redis success:", del1, del2) //del1、del2代表删除了几条数据

		//协程删除数据库中的值
		go func() {
			err := dao.DeleteComment(commentId)
			if err != nil {
				log.Println(err)
			}
		}()
		return nil
	}
	//不在内存中，则直接走数据库删除
	return dao.DeleteComment(commentId)
}

// GetList
// 4、查看评论列表-返回评论list
func (c CommentServiceImpl) GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	log.Println("CommentService-GetList: running") //函数已运行

	//法一、使用SQL语句查询评论列表及用户信息，嵌套user信息。且导致提高耦合性。
	//1.查找CommentData结构体的信息
	commentData := make([]CommentData, 1)
	err := dao.Db.Raw("select T.cid id,T.user_id user_id,T.`name`,T.follow_count,T.follower_count,"+
		"\nif(f.cancel is null,'false','true') is_follow,"+
		"\nT.comment_text content,T.create_date"+
		"\nfrom follows f right join\n("+
		"\n\tselect cid,vid,id user_id,`name`,comment_text,create_date,"+
		"\n\tcount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
		"\n\tcount(if(tag = 'follow' and cancel is not null,1,null)) follow_count"+
		"\n\tfrom\n\t("+
		"\n\t\tselect c.id cid,u.id,c.video_id vid,`name`,f.cancel,comment_text,create_date,'follower' tag"+
		"\n\t\tfrom comments c join users u on c.user_id = u.id and c.cancel = 0"+
		"\n\t\tleft join follows f on u.id = f.user_id and f.cancel = 0"+
		"\n\t\tunion all"+
		"\n\t\tselect c.id cid,u.id,c.video_id vid,`name`,f.cancel,comment_text,create_date,'follow' tag"+
		"\n\t\tfrom comments c join users u on c.user_id = u.id and c.cancel = 0"+
		"\n\t\tleft join follows f on u.id = f.follower_id and f.cancel = 0"+
		"\n\t\t) T\n\t\tgroup by cid,vid,id,`name`,comment_text,create_date"+
		"\n) T on f.follower_id = T.user_id and f.cancel = 0 and f.user_id = ?"+
		"\nwhere vid = ? group by cid order by create_date desc", userId, videoId).Scan(&commentData).Error

	if nil != err {
		log.Println("CommentService-GetList: sql error") //sql查询出错
		return nil, err
	}
	//2.拼接
	commentInfoList := make([]CommentInfo, 0, len(commentData))
	for _, comment := range commentData {
		userData := User{
			Id:            comment.Id,
			Name:          comment.Name,
			FollowCount:   comment.FollowCount,
			FollowerCount: comment.FollowerCount,
			IsFollow:      comment.IsFollow,
		}
		commentData := CommentInfo{
			Id:         comment.Id,
			UserInfo:   userData,
			Content:    comment.Content,
			CreateDate: comment.CreateDate.Format(config.DateTime),
		}
		//3.组装list
		commentInfoList = append(commentInfoList, commentData)
	}

	log.Println("CommentService-GetList: get list success") //成功查询到评论列表

	//查询redis中是否有此记录，无则加入
	//将评论id切片存入redis，
	go func() {
		//评论id循环存入redis，
		for _, _comment := range commentData {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), strconv.Itoa(int(_comment.Id)))
		}
		log.Println("comment list save ids in redis")
	}()

	return commentInfoList, nil

	/*
		//法二：调用dao，先查评论，再循环查用户信息：
		//1.先查询评论列表信息
		commentList, err := dao.GetCommentList(videoId)
		if err != nil {
			log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
			return nil, err
		}
		//提前定义好切片长度
		commentInfoList := make([]CommentInfo, 0, len(commentList))
		for _, comment := range commentList {
			//2.根据查询到的评论用户id和当前用户id，查询评论用户信息
			impl := UserServiceImpl{
				FollowService: &FollowServiceImp{},
			}
			userData, err := impl.GetUserByIdWithCurId(comment.UserId, userId)
			//查看传入评论的两个userid
			//log.Printf("comment.User_id:%v\n", comment.User_id)
			//log.Printf("now_userId:%v\n", userId)
			if err != nil {
				log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
				return nil, err
			}
			commentData := CommentInfo{
				Id:         comment.Id,
				UserInfo:   userData,
				Content:    comment.CommentText,
				CreateDate: comment.CreateDate.Format(config.DateTime),
			}
			//3.组装list
			commentInfoList = append(commentInfoList, commentData)
		}
		log.Println("CommentService-GetList: return list success") //函数执行成功，返回正确信息
		return commentInfoList, nil*/
}

//在redis中存储video_id对应的comment_id 、 comment_id对应的video_id
func insertRedisVideoCommentId(videoId string, commentId string) {
	middleware.RdbVCid.SAdd(middleware.Ctx, videoId, commentId)
	middleware.RdbCVid.Set(middleware.Ctx, commentId, videoId, 0)
}
