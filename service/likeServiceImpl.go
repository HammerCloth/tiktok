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

//IsFavourite 根据userId,videoId查询点赞状态 这边可以快一点,通过查询两个Redis DB;
//step1：查询Redis DB5(key:strUserId)是否已经加载过此信息，通过是否存在value:videoId 判断点赞状态;
//step2:如DB5没有对应信息，查询DB6(key：strVideoId)是否已经加载过此信息，通过是否存在value:userId 判断点赞状态;
//step3:DB5 DB6中都没有对应key,维护DB5对应key，并通过查询key：strUserId中是否存在value:videoId 判断点赞状态;
func (like *LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	//将int64 userId转换为 string strUserId
	strUserId := strconv.FormatInt(userId, 10)
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	//step1:查询Redis DB5,key：strUserId中是否存在value:videoId,key中存在value 返回true，不存在返回false
	if n, err := middleware.Rdb5.Exists(middleware.Ctx, strUserId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("方法:IsFavourite RedisDB5 query key失败：%v", err)
			return false, err
		}
		exist, err1 := middleware.Rdb5.SIsMember(middleware.Ctx, strUserId, videoId).Result()
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err1 != nil {
			log.Printf("方法:IsFavourite RedisDB5 query value失败：%v", err1)
			return false, err1
		}
		log.Printf("方法:IsFavourite RedisDB5 query value成功")
		return exist, nil
	} else { //step2:DB5不存在key,查询Redis DB6,key中存在value 返回true，不存在返回false
		if n, err := middleware.Rdb6.Exists(middleware.Ctx, strVideoId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回默认false,返回错误信息
			if err != nil {
				log.Printf("方法:IsFavourite RedisDB6 query key失败：%v", err)
				return false, err
			}
			exist, err1 := middleware.Rdb6.SIsMember(middleware.Ctx, strVideoId, userId).Result()
			//如果有问题，说明查询redis失败,返回默认false,返回错误信息
			if err1 != nil {
				log.Printf("方法:IsFavourite RedisDB6 query value失败：%v", err1)
				return false, err1
			}
			log.Printf("方法:IsFavourite RedisDB6 query value成功")
			return exist, nil
		} else { //step3:DB5 DB6中都没有对应key，通过userId查询likes表,返回所有点赞videoId，并维护到Redis DB5(key:strUserId)
			videoIdList, err := dao.GetLikeVideoIdList(userId)
			//如果有问题，说明查询数据库失败，返回默认false,返回错误信息："get likeVideoIdList failed"
			if err != nil {
				log.Printf(err.Error())
				return false, err
			}
			//维护Redis DB5(key:strUserId)，遍历videoIdList加入
			for _, likeVideoId := range videoIdList {
				middleware.Rdb5.SAdd(middleware.Ctx, strUserId, likeVideoId)
			}
			//查询Redis DB5,key：strUserId中是否存在value:videoId,存在返回true,不存在返回false
			exist, err1 := middleware.Rdb5.SIsMember(middleware.Ctx, strUserId, videoId).Result()
			//如果有问题，说明操作redis失败,返回默认false,返回错误信息
			if err1 != nil {
				log.Printf("方法:IsFavourite RedisDB5 query value失败：%v", err1)
				return false, err1
			}
			log.Printf("方法:IsFavourite RedisDB5 query value成功")
			return exist, nil
		}
	}
}

//FavouriteCount 根据videoId获取对应点赞数量;
//step1：查询Redis DB6(key:strVideoId)是否已经加载过此信息，通过set集合中userId个数，获取点赞数量;
//step2：DB6中都没有对应key，维护DB6对应key，再通过set集合中userId个数，获取点赞数量;
func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	//step1 如果key:strVideoId存在 则计算集合中userId个数
	if n, err := middleware.Rdb6.Exists(middleware.Ctx, strVideoId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteCount RedisDB6 query key失败：%v", err)
			return 0, err
		}
		//获取集合中userId个数
		count, err1 := middleware.Rdb6.SCard(middleware.Ctx, strVideoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err1 != nil {
			log.Printf("方法:FavouriteCount RedisDB6 query count 失败：%v", err1)
			return 0, err1
		}
		log.Printf("方法:FavouriteCount RedisDB6 query count 成功")
		return count, nil
	} else { //如果Redis DB6不存在此key,通过videoId查询likes表,返回所有点赞userId，并维护到Redis DB6(key:strVideoId)
		//再通过set集合中userId个数,获取点赞数量
		userIdList, err := dao.GetLikeUserIdList(videoId)
		//如果有问题，说明查询数据库失败，返回默认0,返回错误信息："get likeUserIdList failed"
		if err != nil {
			log.Printf(err.Error())
			return 0, err
		}
		//维护Redis DB6(key:strVideoId)，遍历userIdList加入
		for _, likeUserId := range userIdList {
			middleware.Rdb6.SAdd(middleware.Ctx, strVideoId, likeUserId)
		}
		//再通过set集合中userId个数,获取点赞数量
		count, err1 := middleware.Rdb6.SCard(middleware.Ctx, strVideoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err1 != nil {
			log.Printf("方法:FavouriteCount RedisDB6 query count 失败：%v", err1)
			return 0, err
		}
		log.Printf("方法:FavouriteCount RedisDB6 query count 成功")
		return count, nil
	}
}

// FavouriteAction 根据userId，videoId,actionType对视频进行点赞或者取消赞操作;
//step1：更新数据库likes表;
//step2: 维护Redis DB5(key:strUserId),添加或者删除value:videoId,DB6(key:strVideoId),添加或者删除value:userId;
func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	//将int64 videoId转换为 string strVideoId
	strUserId := strconv.FormatInt(userId, 10)
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	//将要操作数据库likes表的信息打入消息队列RmqLikeAdd或者RmqLikeDel
	//拼接打入信息
	sb := strings.Builder{}
	sb.WriteString(strUserId)
	sb.WriteString(" ")
	sb.WriteString(strVideoId)
	//如果是点赞操作，消息打入RmqLikeAdd队列
	if actionType == config.LikeAction {
		middleware.RmqLikeAdd.Publish(sb.String())
	} else { //如果是取消赞操作，消息打入RmqLikeDel队列
		middleware.RmqLikeDel.Publish(sb.String())
	}

	//step2:维护Redis DB5、DB6;
	//执行点赞操作维护
	if actionType == config.LikeAction {
		//查询Redis DB5(key:strUserId)是否已经加载过此信息
		if n, err := middleware.Rdb5.Exists(middleware.Ctx, strUserId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回默认false,返回错误信息
			if err != nil {
				log.Printf("方法:FavouriteAction RedisDB5 query key失败：%v", err)
				return err
			} //
			middleware.Rdb5.SAdd(middleware.Ctx, strUserId, videoId)
		} else { //如果不存在，则新建key，加入点赞videoid,搜索数据库videoid，然后依次加入
			middleware.Rdb5.SAdd(middleware.Ctx, strUserId, videoId)
			videoIdList, err := dao.GetLikeVideoIdList(userId)
			if err != nil {
				return err
			}
			for _, likeVideoId := range videoIdList {
				middleware.Rdb5.SAdd(middleware.Ctx, strUserId, likeVideoId)
			}
		}
		//step2  如果videoId存在 则加入点赞userid
		if n, _ := middleware.Rdb6.Exists(middleware.Ctx, strVideoId).Result(); n > 0 {
			middleware.Rdb6.SAdd(middleware.Ctx, strVideoId, userId)
		} else { //如果不存在，则新建key，加入点赞userid,搜索数据库userid，然后依次加入
			middleware.Rdb6.SAdd(middleware.Ctx, strVideoId, userId)
			userIdList, err := dao.GetLikeUserIdList(videoId)
			if err != nil {
				return err
			}
			for _, likeUserId := range userIdList {
				middleware.Rdb6.SAdd(middleware.Ctx, strVideoId, likeUserId)
			}
		}
	} else { //取消赞
		//step1  如果userid存在 则删除点赞videoid
		if n, _ := middleware.Rdb5.Exists(middleware.Ctx, strUserId).Result(); n > 0 {
			middleware.Rdb5.SRem(middleware.Ctx, strUserId, videoId)
		} else { //如果不存在，则新建key,搜索数据库videoid，然后依次加入,再删除点赞videoid
			videoIdList, err := dao.GetLikeVideoIdList(userId)
			if err != nil {
				return err
			}
			for _, likeVideoId := range videoIdList {
				middleware.Rdb5.SAdd(middleware.Ctx, strUserId, likeVideoId)
			}
			middleware.Rdb5.SRem(middleware.Ctx, strUserId, videoId)
		}
		//step2  如果videoId存在 则删除点赞userid
		if n, _ := middleware.Rdb6.Exists(middleware.Ctx, strVideoId).Result(); n > 0 {
			middleware.Rdb6.SRem(middleware.Ctx, strVideoId, userId)
		} else { //如果不存在，则新建key,搜索数据库userid，然后依次加入,再删除点赞userid,
			userIdList, err := dao.GetLikeUserIdList(videoId)
			if err != nil {
				return err
			}
			for _, likeUserId := range userIdList {
				middleware.Rdb6.SAdd(middleware.Ctx, strVideoId, likeUserId)
			}
			middleware.Rdb6.SRem(middleware.Ctx, strVideoId, userId)
		}
	}
	return nil
}

//GetFavouriteList 根据userId，curId(当前用户Id),返回userId的点赞列表;
//step1：查询Redis DB5(key:strUserId)是否已经加载过此信息，获取集合中全部videoId，并添加到点赞列表集合中;
//step2：DB5中都没有对应key，维护DB5对应key，同时添加到点赞列表集合中;
func (like *LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	//将int64 userId转换为 string strUserId
	strUserId := strconv.FormatInt(userId, 10)
	//step1:查询Redis DB5,如果key：strUserId存在,则获取集合中全部videoId
	if n, err := middleware.Rdb5.Exists(middleware.Ctx, strUserId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err != nil {
			log.Printf("方法:GetFavouriteList RedisDB6 query key失败：%v", err)
			return nil, err
		}
		//获取集合中全部videoId
		videoIdList, err1 := middleware.Rdb5.SMembers(middleware.Ctx, strUserId).Result()
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err1 != nil {
			log.Printf("方法:GetFavouriteList RedisDB6 get values失败：%v", err1)
			return nil, err1
		}
		//提前定义好切片长度,生成点赞列表集合
		favoriteVideoList := make([]Video, 0, len(videoIdList))
		//如果查询成功，无论是否有数据，遍历 string videoIdList,获得其中的 string videoId；
		for _, likeVideoId := range videoIdList {
			//将string likeVideoId转换为 int64 VideoId
			VideoId, _ := strconv.ParseInt(likeVideoId, 10, 64)
			//调用videoService接口，GetVideo：根据videoId，当前用户id:curId，返回Video类型对象
			video, err2 := like.GetVideo(VideoId, curId)
			//如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且跳过此视频
			if err2 != nil {
				log.Println(errors.New("this favourite video is miss"))
				continue
			}
			//将每个Video类型对象添加到集合中去
			favoriteVideoList = append(favoriteVideoList, video)
		}
		return favoriteVideoList, nil
	} else { //如果Redis DB5不存在此key,通过userId查询likes表,返回所有点赞videoId，并维护到Redis DB5(key:strUserId)
		videoIdList, err := dao.GetLikeVideoIdList(userId)
		//如果有问题，说明查询数据库失败，返回nil和错误信息:"get likeVideoIdList failed"
		if err != nil {
			log.Println(err.Error())
			return nil, err
		}
		//提前定义好切片长度,生成点赞列表集合
		favoriteVideoList := make([]Video, 0, len(videoIdList))
		//如果查询成功，无论是否有数据,遍历 int videoIdList,获得其中的 int likeVideoId，维护到Redis DB5(key:strUserId)
		//同时添加到点赞列表集合中
		for _, likeVideoId := range videoIdList {
			middleware.Rdb5.SAdd(middleware.Ctx, strUserId, likeVideoId)
			video, err1 := like.GetVideo(likeVideoId, curId)
			//如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且跳过此视频
			if err1 != nil {
				log.Println(errors.New("can't find this favourite video"))
				continue
			}
			//将每个Video类型对象添加到集合中去
			favoriteVideoList = append(favoriteVideoList, video)
		}
		return favoriteVideoList, nil
	}
}
