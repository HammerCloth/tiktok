package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware"
	"errors"
	"log"
	"strconv"
	"strings"
)

type LikeServiceImpl struct {
	VideoService
}

// 根据userid,videoid查询点赞信息 这边可以快一点（修改下）判断两个redis
//Redis是否存在userId set{videoid} 点赞的视频id 功能： 判断是否存在该videoid来判断是否点赞
func (like *LikeServiceImpl) IsFavourit(videoId int64, userId int64) (bool, error) {
	suserId := strconv.FormatInt(userId, 10)
	svideoId := strconv.FormatInt(videoId, 10)
	//step1  如果userid存在 则判断videoid是否存在
	if n, _ := middleware.Rdb5.Exists(middleware.Ctx, suserId).Result(); n > 0 {
		exist, err := middleware.Rdb5.SIsMember(middleware.Ctx, suserId, videoId).Result()
		//如果有问题，说明操作redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("RedisIsFavourit(videoId) 失败：%v", err)
			return false, err
		}
		log.Printf("RedisIsFavourit(videoId) 成功")
		return exist, nil
	} else { //如果不存在，则判断redis videoid 是否存在
		if n, _ := middleware.Rdb6.Exists(middleware.Ctx, svideoId).Result(); n > 0 {
			exist, err := middleware.Rdb6.SIsMember(middleware.Ctx, svideoId, userId).Result()
			//如果有问题，说明操作redis失败,返回默认false,返回错误信息
			if err != nil {
				log.Printf("RedisIsFavourit(videoId) 失败：%v", err)
				return false, err
			}
			log.Printf("RedisIsFavourit(videoId) 成功")
			return exist, nil
		} else { // 还不存在则新建key，搜索数据库videoid，然后依次加入，再判断videoid是否存在
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			if err1 != nil {
				log.Printf(err1.Error())
				return false, err1
			}
			for _, likeVideoId := range videoIdList {
				middleware.Rdb5.SAdd(middleware.Ctx, suserId, likeVideoId)
			}
			exist, err := middleware.Rdb5.SIsMember(middleware.Ctx, suserId, videoId).Result()
			//如果有问题，说明操作redis失败,返回默认false,返回错误信息
			if err != nil {
				log.Printf("RedisIsFavourit(videoId) 失败：%v", err)
				return false, err
			}
			log.Printf("RedisIsFavourit(videoId) 成功")
			return exist, nil
		}
	}
}

//根据videoid获取点赞数量  Redis是否存在videoId set{userid} 点赞的用户id   功能：计算size得到该视频点赞数
func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	//查询videoid对应点赞数量
	svideoId := strconv.FormatInt(videoId, 10)
	//step1 如果videoId存在 则计算点赞userid
	if n, _ := middleware.Rdb6.Exists(middleware.Ctx, svideoId).Result(); n > 0 {
		count, err := middleware.Rdb6.SCard(middleware.Ctx, svideoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err != nil {
			log.Printf("RedisLikeCount(videoId) 失败：%v", err)
			return 0, err
		}
		log.Printf("RedisLikeCount(videoId) 成功")
		return count, nil
	} else { //如果不存在，则新建key，搜索数据库userid，然后依次加入,然后计算size
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		if err1 != nil {
			log.Printf(err1.Error())
			return 0, err1
		}
		for _, likeUserId := range userIdList {
			middleware.Rdb6.SAdd(middleware.Ctx, svideoId, likeUserId)
		}
		count, err := middleware.Rdb6.SCard(middleware.Ctx, svideoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err != nil {
			log.Printf("RedisLikeCount(videoId) 失败：%v", err)
			return 0, err
		}
		log.Printf("RedisLikeCount(videoId) 成功")
		return count, nil
	}
}

func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, action_type int32) error {
	suserId := strconv.FormatInt(userId, 10)
	svideoId := strconv.FormatInt(videoId, 10)
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(suserId)
	sb.WriteString(" ")
	sb.WriteString(svideoId)
	if action_type == config.Likeaction {
		middleware.RmqLikeAdd.Publish(sb.String())
	} else {
		middleware.RmqLikeDel.Publish(sb.String())
	}
	// 更新redis信息。
	/*
		1-Redis是否存在userId set{videoid} 点赞的视频id   功能： 计算size得到用户的喜欢数，用户的喜欢列表
		2-Redis是否存在videoId set{userid} 点赞的用户id   功能：计算size得到该视频点赞数
	*/

	if action_type == config.Likeaction {
		//step1  如果userid存在 则加入点赞videoid
		if n, _ := middleware.Rdb5.Exists(middleware.Ctx, suserId).Result(); n > 0 {
			middleware.Rdb5.SAdd(middleware.Ctx, suserId, videoId)
		} else { //如果不存在，则新建key，加入点赞videoid,搜索数据库videoid，然后依次加入
			middleware.Rdb5.SAdd(middleware.Ctx, suserId, videoId)
			videoIdList, err := dao.GetLikeVideoIdList(userId)
			if err != nil {
				return err
			}
			for _, likeVideoId := range videoIdList {
				middleware.Rdb5.SAdd(middleware.Ctx, suserId, likeVideoId)
			}
		}
		//step2  如果videoId存在 则加入点赞userid
		if n, _ := middleware.Rdb6.Exists(middleware.Ctx, svideoId).Result(); n > 0 {
			middleware.Rdb6.SAdd(middleware.Ctx, svideoId, userId)
		} else { //如果不存在，则新建key，加入点赞userid,搜索数据库userid，然后依次加入
			middleware.Rdb6.SAdd(middleware.Ctx, svideoId, userId)
			userIdList, err := dao.GetLikeUserIdList(videoId)
			if err != nil {
				return err
			}
			for _, likeUserId := range userIdList {
				middleware.Rdb6.SAdd(middleware.Ctx, svideoId, likeUserId)
			}
		}
	} else { //取消赞
		//step1  如果userid存在 则删除点赞videoid
		if n, _ := middleware.Rdb5.Exists(middleware.Ctx, suserId).Result(); n > 0 {
			middleware.Rdb5.SRem(middleware.Ctx, suserId, videoId)
		} else { //如果不存在，则新建key,搜索数据库videoid，然后依次加入,再删除点赞videoid
			videoIdList, err := dao.GetLikeVideoIdList(userId)
			if err != nil {
				return err
			}
			for _, likeVideoId := range videoIdList {
				middleware.Rdb5.SAdd(middleware.Ctx, suserId, likeVideoId)
			}
			middleware.Rdb5.SRem(middleware.Ctx, suserId, videoId)
		}
		//step2  如果videoId存在 则删除点赞userid
		if n, _ := middleware.Rdb6.Exists(middleware.Ctx, svideoId).Result(); n > 0 {
			middleware.Rdb6.SRem(middleware.Ctx, svideoId, userId)
		} else { //如果不存在，则新建key,搜索数据库userid，然后依次加入,再删除点赞userid,
			userIdList, err := dao.GetLikeUserIdList(videoId)
			if err != nil {
				return err
			}
			for _, likeUserId := range userIdList {
				middleware.Rdb6.SAdd(middleware.Ctx, svideoId, likeUserId)
			}
			middleware.Rdb6.SRem(middleware.Ctx, svideoId, userId)
		}
	}
	return nil
}

func (like *LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	//1.先查询点赞列表信息  Redis是否存在userId set{videoid} 点赞的视频id
	//功能： 计算size得到用户的喜欢数，用户的喜欢列表
	suserId := strconv.FormatInt(userId, 10)
	//step1  如果userid存在 则获取集合中全部videoid
	if n, _ := middleware.Rdb5.Exists(middleware.Ctx, suserId).Result(); n > 0 {
		videoIdList, err := middleware.Rdb5.SMembers(middleware.Ctx, suserId).Result()
		if err != nil {
			log.Printf("RedisGetFavouriteList(userId) 失败：%v", err)
			return nil, err
		}
		//提前定义好切片长度,生成集合
		favorite_videolist := make([]Video, 0, len(videoIdList))
		//如果查询成功，无论是否有数据，遍历likelist,获得其中的video_id；
		//测试结构体，协同开发
		//likesub := new(LikeSub)
		for _, likeVideoId := range videoIdList {
			//测试函数，协同开发
			//video, err1 := likesub.GetVideo(likedata.Video_id,userId)
			//调用video接口，Getvideo：根据videoid，当前用户id，返回video对象
			VideoId, _ := strconv.ParseInt(likeVideoId, 10, 64)
			video, err1 := like.GetVideo(VideoId, curId)
			if err1 != nil { //如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且跳过
				log.Println(errors.New("can't find this favourite video"))
				continue
			} //将每个video对象添加到集合中去
			favorite_videolist = append(favorite_videolist, video)
		}
		return favorite_videolist, nil
	} else { //如果不存在，则新建key，搜索数据库videoid，得到videoIdList，然后依次加入
		videoIdList, err := dao.GetLikeVideoIdList(userId)
		//如果有问题，说明查询数据库失败，返回空和错误err:"get likeList failed"
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		for _, likeVideoId := range videoIdList {
			middleware.Rdb5.SAdd(middleware.Ctx, suserId, likeVideoId)
			//log.Printf("video_id：%t,%v", string(likeVideoId), string(likeVideoId))
		}
		//提前定义好切片长度,生成集合
		favorite_videolist := make([]Video, 0, len(videoIdList))
		//如果查询成功，无论是否有数据，遍历likelist,获得其中的video_id；
		//测试结构体，协同开发
		//likesub := new(LikeSub)
		for _, likeVideoId := range videoIdList {
			//测试函数，协同开发
			//video, err1 := likesub.GetVideo(likedata.Video_id,userId)
			//调用video接口，Getvideo：根据videoid，当前用户id，返回video对象
			video, err1 := like.GetVideo(likeVideoId, curId)
			if err1 != nil { //如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且跳过
				log.Println(errors.New("can't find this favourite video"))
				continue
			} //将每个video对象添加到集合中去
			favorite_videolist = append(favorite_videolist, video)
		}
		return favorite_videolist, nil
	}
}
