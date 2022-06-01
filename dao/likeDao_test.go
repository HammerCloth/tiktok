package dao

import (
	"fmt"
	"testing"
)

func TestGetLikeUserIdList(t *testing.T) {
	Init()
	list, err := GetLikeUserIdList(54)
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestUpdateLike(t *testing.T) {
	Init()
	err := UpdateLike(3, 54, 0)
	fmt.Printf("%v", err)
}

func TestInsertLike(t *testing.T) {
	Init()
	err := InsertLike(Like{
		UserId:  20003,
		VideoId: 71,
		Cancel:  0,
	})
	fmt.Printf("%v", err)
}

func TestGetLikeInfo(t *testing.T) {
	Init()
	likeInfo, err := GetLikeInfo(3, 71)
	fmt.Printf("%v", likeInfo)
	fmt.Printf("%v", err)
}

func TestGetLikeVideoIdList(t *testing.T) {
	Init()
	videoIdList, err := GetLikeVideoIdList(3)
	fmt.Printf("%v", videoIdList)
	fmt.Printf("%v", err)
}
