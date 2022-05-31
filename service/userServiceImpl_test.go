package service

import (
	"TikTok/dao"
	"fmt"
	"testing"
)

func TestGetTableUserList(t *testing.T) {
	impl := UserServiceImpl{}
	list := impl.GetTableUserList()
	fmt.Printf("%v", list)
}

func TestGetTableUserByUsername(t *testing.T) {
	impl := UserServiceImpl{}
	list := impl.GetTableUserByUsername("test")
	fmt.Printf("%v", list)
}

func TestGetTableUserById(t *testing.T) {
	impl := UserServiceImpl{}
	list := impl.GetTableUserById(int64(4))
	fmt.Printf("%v", list)
}

func TestInsertTableUser(t *testing.T) {
	impl := UserServiceImpl{}
	user := &dao.TableUser{
		Id:       20000,
		Name:     "qaq",
		Password: "111111",
	}
	list := impl.InsertTableUser(user)
	fmt.Printf("%v", list)
}

func TestGetUserById(t *testing.T) {
	impl := UserServiceImpl{
		FollowService: &FollowServiceImp{},
		LikeService:   &LikeServiceImpl{},
	}
	list, _ := impl.GetUserById(int64(4))
	fmt.Printf("%v", list)
}

func TestGetUserByIdWithCurId(t *testing.T) {
	impl := UserServiceImpl{
		FollowService: &FollowServiceImp{},
	}
	list, _ := impl.GetUserByIdWithCurId(int64(482), int64(130))
	fmt.Printf("%v", list)
}
