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
	likedata := new(dao.Like)
	//未查询到数据，返回未点赞；
	if result := dao.Db.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).First(&likedata); result.RowsAffected == 0 {
		return false, errors.New("can't find this data")
	} //查询到数据，根据Cancel值判断是否点赞；
	if likedata.Cancel == config.Islike {
		return true, nil
	} else {
		return false, nil
	}
}

func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	var count int64
	err := dao.Db.Table("likes").Where("video_id = ? AND cancel = ?", videoId, config.Islike).
		Count(&count).Error
	//当查询出现异常时
	if err != nil {
		return 0, errors.New("An unknown exception occurred in the query")
	} else {
		return count, nil
	}
}

func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, action_type int32) error {
	likedata := new(dao.Like)
	//先查询是否有这条数据。
	result := dao.Db.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).First(&likedata)
	//点赞行为,否则取消赞
	if action_type == config.Likeaction {
		//没查到这条数据，则新建这条数据,否则更新即可;
		if result.RowsAffected == 0 {
			likedata.User_id = userId
			likedata.Video_id = videoId
			likedata.Cancel = config.Islike
			if result1 := dao.Db.Table("likes").Create(&likedata); result1.RowsAffected == 0 {
				return errors.New("insert data fail")
			}
		} else {
			if result2 := dao.Db.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).
				Update("cancel", config.Islike); result2.RowsAffected == 0 {
				return errors.New("update data fail")
			}
		}
	} else {
		//只有当前是点赞状态才能取消点赞这个行为，如果查询不到数据则返回错误；
		if result.RowsAffected == 0 {
			return errors.New("can't find this data")
		} else {
			if result3 := dao.Db.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).
				Update("cancel", config.Unlike); result3.RowsAffected == 0 {
				return errors.New("update data fail")
			}
		}
	}
	return nil
}

func (like *LikeServiceImpl) GetFavouriteList(userId int64) ([]Video, error) {
	var favorite_videolist []Video
	var video_ids []int64
	if result := dao.Db.Table("likes").Select("video_id").Where("user_id = ? AND cancel = ?", userId, config.Islike).
		Find(&video_ids); result.RowsAffected == 0 { //如果查询不到数据，说明这个用户没有点赞视频，返回空；
		return favorite_videolist, errors.New("can't find favourite video")
	}
	//如果查询到数据，遍历video_id,将每个video对象添加到集合中去；
	//likesub := new(LikeSub)
	for video_id := range video_ids {
		//video, err1 := likesub.GetVideo(int64(video_id),userId)
		video, err1 := like.GetVideo(int64(video_id), userId)
		if err1 != nil { //如果没有获取这个video_id的视频，视频可能被删除了,跳过
			log.Panicln(errors.New("can't find this favourite video"))
			continue
		}
		favorite_videolist = append(favorite_videolist, video)
	}
	return favorite_videolist, nil
}
