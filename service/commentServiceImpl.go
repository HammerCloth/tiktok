package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/rabbitmq"
	"TikTok/middleware/redis"
	"log"
	"sort"
	"strconv"
	"sync"
	"time"
)

type CommentServiceImpl struct {
	UserService
}

// CountFromVideoId
// 1、使用video id 查询Comment数量
func (c CommentServiceImpl) CountFromVideoId(videoId int64) (int64, error) {
	//先在缓存中查
	cnt, err := redis.RdbVCid.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
	if err != nil { //若查询缓存出错，则打印log
		//return 0, err
		log.Println("count from redis error:", err)
	}
	log.Println("comment count redis :", cnt)
	//1.缓存中查到了数量，则返回数量值-1（去除0值）
	if cnt != 0 {
		return cnt - 1, nil
	}
	//2.缓存中查不到则去数据库查
	cntDao, err1 := dao.Count(videoId)
	log.Println("comment count dao :", cntDao)
	if err1 != nil {
		log.Println("comment count dao err:", err1)
		return 0, nil
	}
	//将评论id切片存入redis-第一次存储 V-C set 值：
	go func() {
		//查询评论id list
		cList, _ := dao.CommentIdList(videoId)
		//先在redis中存储一个-1值，防止脏读
		_, _err := redis.RdbVCid.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), config.DefaultRedisValue).Result()
		if _err != nil { //若存储redis失败，则直接返回
			log.Println("redis save one vId - cId 0 failed")
			return
		}
		//设置key值过期时间
		_, err := redis.RdbVCid.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		//评论id循环存入redis
		for _, commentId := range cList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), commentId)
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
	commentInfo.CreateDate = comment.CreateDate   //评论时间

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
	//将此发表的评论id存入redis
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
	//1.先查询redis，若有则删除，返回客户端-再go协程删除数据库；无则在数据库中删除，返回客户端。
	n, err := redis.RdbCVid.Exists(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
	if err != nil {
		log.Println(err)
	}
	if n > 0 { //在缓存中有此值，则找出来删除，然后返回
		vid, err1 := redis.RdbCVid.Get(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err1 != nil { //没找到，返回err
			log.Println("redis find CV err:", err1)
		}
		//删除，两个redis都要删除
		del1, err2 := redis.RdbCVid.Del(redis.Ctx, strconv.FormatInt(commentId, 10)).Result()
		if err2 != nil {
			log.Println(err2)
		}
		del2, err3 := redis.RdbVCid.SRem(redis.Ctx, vid, strconv.FormatInt(commentId, 10)).Result()
		if err3 != nil {
			log.Println(err3)
		}
		log.Println("del comment in Redis success:", del1, del2) //del1、del2代表删除了几条数据

		/*
			//由于多协程会造成数据库读取压力大，因此使用chan来减少压力
			commentChan := make(chan int64)
			//协程删除数据库中的值
			go func() {
				commentChan <- commentId
				close(commentChan)
				//err := dao.DeleteComment(commentId)
				//if err != nil {
				//	log.Println(err)
				//}
			}()
			go func() {
				for cId := range commentChan {
					err := dao.DeleteComment(cId)
					if err != nil {
						log.Println(err)
					}
				}
			}()*/

		//使用mq进行数据库中评论的删除-评论状态更新
		//评论id传入消息队列
		rabbitmq.RmqCommentDel.Publish(strconv.FormatInt(commentId, 10))
		return nil
	}
	//不在内存中，则直接走数据库删除
	return dao.DeleteComment(commentId)
}

// GetList
// 4、查看评论列表-返回评论list
func (c CommentServiceImpl) GetList(videoId int64, userId int64) ([]CommentInfo, error) {
	log.Println("CommentService-GetList: running") //函数已运行
	/*
		//法一、使用SQL语句查询评论列表及用户信息，嵌套user信息。且导致提高耦合性。
		//1.查找CommentData结构体的信息
		commentData := make([]CommentData, 1)
		result := dao.Db.Raw("select T.cid id,T.user_id user_id,T.`name`,T.follow_count,T.follower_count,"+
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
			"\nwhere vid = ? group by cid order by create_date desc", userId, videoId).Scan(&commentData)

		err := result.Error

		if nil != err {
			log.Println("CommentService-GetList: sql error") //sql查询出错
			return nil, err
		}
		//当前有0条评论
		if result.RowsAffected == 0 {
			return nil, nil
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
			_commentInfo := CommentInfo{
				Id:         comment.Id,
				UserInfo:   userData,
				Content:    comment.Content,
				CreateDate: comment.CreateDate.Format(config.DateTime),
			}
			//3.组装list
			commentInfoList = append(commentInfoList, _commentInfo)
		}
		//-----------------------法一结束--------------------------
	*/

	//法二：调用dao，先查评论，再循环查用户信息：
	//1.先查询评论列表信息
	commentList, err := dao.GetCommentList(videoId)
	if err != nil {
		log.Println("CommentService-GetList: return err: " + err.Error()) //函数返回提示错误信息
		return nil, err
	}
	//当前有0条评论
	if commentList == nil {
		return nil, nil
	}

	//提前定义好切片长度
	commentInfoList := make([]CommentInfo, len(commentList))

	wg := &sync.WaitGroup{}
	wg.Add(len(commentList))
	idx := 0
	for _, comment := range commentList {
		//2.调用方法组装评论信息，再append
		var commentData CommentInfo
		//将评论信息进行组装，添加想要的信息,插入从数据库中查到的数据
		go func(comment dao.Comment) {
			oneComment(&commentData, &comment, userId)
			//3.组装list
			//commentInfoList = append(commentInfoList, commentData)
			commentInfoList[idx] = commentData
			idx = idx + 1
			wg.Done()
		}(comment)
	}
	wg.Wait()
	//评论排序-按照主键排序
	sort.Sort(CommentSlice(commentInfoList))
	//------------------------法二结束----------------------------

	//协程查询redis中是否有此记录，无则将评论id切片存入redis
	go func() {
		//1.先在缓存中查此视频是否已有评论列表
		cnt, err1 := redis.RdbVCid.SCard(redis.Ctx, strconv.FormatInt(videoId, 10)).Result()
		if err1 != nil { //若查询缓存出错，则打印log
			//return 0, err
			log.Println("count from redis error:", err)
		}
		//2.缓存中查到了数量大于0，则说明数据正常，不用更新缓存
		if cnt > 0 {
			return
		}
		//3.缓存中数据不正确，更新缓存：
		//先在redis中存储一个-1 值，防止脏读
		_, _err := redis.RdbVCid.SAdd(redis.Ctx, strconv.Itoa(int(videoId)), config.DefaultRedisValue).Result()
		if _err != nil { //若存储redis失败，则直接返回
			log.Println("redis save one vId - cId 0 failed")
			return
		}
		//设置key值过期时间
		_, err2 := redis.RdbVCid.Expire(redis.Ctx, strconv.Itoa(int(videoId)),
			time.Duration(config.OneMonth)*time.Second).Result()
		if err2 != nil {
			log.Println("redis save one vId - cId expire failed")
		}
		//将评论id循环存入redis
		for _, _comment := range commentInfoList {
			insertRedisVideoCommentId(strconv.Itoa(int(videoId)), strconv.Itoa(int(_comment.Id)))
		}
		log.Println("comment list save ids in redis")
	}()

	log.Println("CommentService-GetList: return list success") //函数执行成功，返回正确信息
	return commentInfoList, nil
}

//在redis中存储video_id对应的comment_id 、 comment_id对应的video_id
func insertRedisVideoCommentId(videoId string, commentId string) {
	//在redis-RdbVCid中存储video_id对应的comment_id
	_, err := redis.RdbVCid.SAdd(redis.Ctx, videoId, commentId).Result()
	if err != nil { //若存储redis失败-有err，则直接删除key
		log.Println("redis save send: vId - cId failed, key deleted")
		redis.RdbVCid.Del(redis.Ctx, videoId)
		return
	}
	//在redis-RdbCVid中存储comment_id对应的video_id
	_, err = redis.RdbCVid.Set(redis.Ctx, commentId, videoId, 0).Result()
	if err != nil {
		log.Println("redis save one cId - vId failed")
	}
}

//此函数用于给一个评论赋值：评论信息+用户信息 填充
func oneComment(comment *CommentInfo, com *dao.Comment, userId int64) {
	var wg sync.WaitGroup
	wg.Add(1)
	//根据评论用户id和当前用户id，查询评论用户信息
	impl := UserServiceImpl{
		FollowService: &FollowServiceImp{},
	}
	var err error
	comment.Id = com.Id
	comment.Content = com.CommentText
	comment.CreateDate = com.CreateDate.Format(config.DateTime)
	comment.UserInfo, err = impl.GetUserByIdWithCurId(com.UserId, userId)
	if err != nil {
		log.Println("CommentService-GetList: GetUserByIdWithCurId return err: " + err.Error()) //函数返回提示错误信息
	}
	wg.Done()
	wg.Wait()
}

// CommentSlice 此变量以及以下三个函数都是做排序-准备工作
type CommentSlice []CommentInfo

func (a CommentSlice) Len() int { //重写Len()方法
	return len(a)
}
func (a CommentSlice) Swap(i, j int) { //重写Swap()方法
	a[i], a[j] = a[j], a[i]
}
func (a CommentSlice) Less(i, j int) bool { //重写Less()方法
	return a[i].Id > a[j].Id
}
