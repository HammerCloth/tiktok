package service

import "TikTok/dao"

type VideoSub struct {
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
