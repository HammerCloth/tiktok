package service

import (
	"TikTok/config"
	"TikTok/dao"
	"errors"
	"log"
)

type LikeServiceImpl struct {
	VideoService
}

func (like *LikeServiceImpl) IsFavourit(videoId int64, userId int64) (bool, error) {
	//未查询到数据，返回未点赞；
	likedata, err := dao.NewLikeDaoInstance().GetLikeInfo(userId, videoId)
	if err != nil {
		return false, err
	} else { //查询到数据，根据Cancel值判断是否点赞；
		if likedata.Cancel == config.Islike {
			return true, nil
		} else {
			return false, nil
		}
	}
}

func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	return dao.NewLikeDaoInstance().GetLikeCount(videoId)
}

func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, action_type int32) error {
	var likedata dao.Like
	//先查询是否有这条数据。
	_, err := dao.NewLikeDaoInstance().GetLikeInfo(userId, videoId)
	//点赞行为,否则取消赞
	if action_type == config.Likeaction {
		//没查到这条数据，则新建这条数据,否则更新即可;
		if err != nil {
			likedata.User_id = userId
			likedata.Video_id = videoId
			likedata.Cancel = config.Islike
			return dao.NewLikeDaoInstance().InsertLike(likedata)
		} else {
			return dao.NewLikeDaoInstance().UpdateLike(userId, videoId, config.Islike)
		}
	} else {
		//只有当前是点赞状态才能取消点赞这个行为，如果查询不到数据则返回错误；
		if err != nil {
			log.Println(err.Error())
			return errors.New("can't find this data")
		} else {
			return dao.NewLikeDaoInstance().UpdateLike(userId, videoId, config.Unlike)
		}
	}
	return nil
}

func (like *LikeServiceImpl) GetFavouriteList(userId int64) ([]Video, error) {
	//1.先查询点赞列表信息
	likeList, err := dao.NewLikeDaoInstance().GetLikeList(userId)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	//提前定义好切片长度
	favorite_videolist := make([]Video, 0, len(likeList))
	//如果查询到数据，遍历video_id；
	//likesub := new(LikeSub)
	for _, likedata := range likeList {
		//video, err1 := likesub.GetVideo(likedata.Video_id,userId)
		video, err1 := like.GetVideo(likedata.Video_id, userId)
		if err1 != nil { //如果没有获取这个video_id的视频，视频可能被删除了,跳过
			log.Panicln(errors.New("can't find this favourite video"))
			continue
		} //将每个video对象添加到集合中去
		favorite_videolist = append(favorite_videolist, video)
	}
	return favorite_videolist, nil
}
