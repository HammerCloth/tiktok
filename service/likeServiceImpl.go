package service

import (
	"TikTok/dao"
	"errors"
)

type LikeServiceImpl struct {
	VideoService
}

func (like *LikeServiceImpl) IsFavourit(videoId int64, userId int64) (bool, error) {
	likedata := new(dao.Like)
	if err := dao.Db.Table("likes").Where("user_id = ? AND video_id = ?", userId, videoId).First(&likedata); err != nil {
		return true, errors.New("can't find this data")
	}
	if likedata.Cancel == 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	var count int64
	if err := dao.Db.Table("likes").Where("video_id = ?", videoId).First(&count); err != nil {
		return 0, errors.New("can't find this data")
	}
	return count, nil
}

//func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, action_type int32) error {
//	panic("implement me")
//}
//
//func (like *LikeServiceImpl) GetFavouriteList(userId int64) ([]Video, error) {
//	panic("implement me")
//}
