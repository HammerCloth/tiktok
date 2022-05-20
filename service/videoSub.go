package service

import (
	"TikTok/dao"
)

type VideoSub struct {
}

func (vs VideoSub) Send(comment dao.Comment) error {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) CountFromVideoId(id int64) (int64, error) {
	return 4, nil
}

func (vs VideoSub) DelComment(id int64) error {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetList(vedioId int64, userId int64) ([]CommentInfo, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) IsFollowing(userId int64, targetId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetFollowerCnt(userId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetFollowingCnt(userId int64) (int64, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetFollowing(userId int64) ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetFollowers(userId int64) ([]User, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) CheckCommentString() string {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) IsFavourit(videoId int64, userId int64) (bool, error) {
	return true, nil
}

func (vs VideoSub) FavouriteCount(videoId int64) (int64, error) {
	return 3, nil
}

func (vs VideoSub) FavouriteAction(userId int64, videoId int64, action_type int32) error {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetFavouriteList(userId int64) ([]Video, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetTableUserList() []dao.TableUser {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetTableUserByUsername(name string) dao.TableUser {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetTableUserById(id int64) dao.TableUser {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) InsertTableUser(tableUser *dao.TableUser) bool {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetUserById(id int64) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (vs VideoSub) GetUserByIdWithCurId(id int64, curId int64) (User, error) {
	var user User
	return user, nil
}
