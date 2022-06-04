package dao

import (
	"fmt"
	"testing"
	"time"
)

//1、使用video id 查询Comment数量 的测试函数
func TestCountComment(t *testing.T) {
	Init()
	count, err := Count(14)
	fmt.Printf("%v\n", count)
	fmt.Printf("[%v]", err)
}

//2、发表评论 的测试函数
func TestInsertComment(t *testing.T) {
	Init()
	nowTime := time.Now() //评论时间记录
	comment := Comment{
		UserId:      20008,
		VideoId:     1,
		CommentText: "user20008commentVideo1-2",
		CreateDate:  nowTime,
		Cancel:      0,
	}
	comList, err := InsertComment(comment)
	fmt.Printf("[comList:%v][err:%v]", comList, err)
}

//3、删除评论 的测试函数
func TestDelComment(t *testing.T) {
	Init()
	err := DeleteComment(int64(8))
	fmt.Printf("[%v]", err)
}

//4.根据视频id查询所属评论全部列表信息 的测试函数
func TestCommentList(t *testing.T) {
	Init()
	list, err := GetCommentList(int64(1))
	fmt.Printf("%v\n", list)
	fmt.Printf("[%v]", err)
}
