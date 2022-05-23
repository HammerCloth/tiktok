package dao

import (
	"TikTok/config"
	"fmt"
	"testing"
	"time"
)

//1、使用video id 查询Comment数量 的测试函数
func TestCountComment(t *testing.T) {
	Init()
	count, err := NewCommentDaoInstance().Count(14)
	fmt.Printf("%v\n", count)
	fmt.Printf("[%v]", err)
}

//2、发表评论 的测试函数
func TestInsertComment(t *testing.T) {
	Init()
	nowTime := time.Now().Format(config.DateTime) //评论时间记录
	comment := Comment{
		User_id:      20008,
		Video_id:     1,
		Comment_text: "user20008commentVideo1-2",
		Create_date:  nowTime,
		Cancel:       0,
	}
	err := NewCommentDaoInstance().InsertComment(comment)
	fmt.Printf("[%v]", err)
}

//3、删除评论 的测试函数
func TestDelComment(t *testing.T) {
	Init()
	err := NewCommentDaoInstance().DeleteComment(int64(8))
	fmt.Printf("[%v]", err)
}

//4.根据视频id查询所属评论全部列表信息 的测试函数
func TestCommentList(t *testing.T) {
	Init()
	list, err := NewCommentDaoInstance().GetCommentList(int64(1))
	fmt.Printf("%v\n", list)
	fmt.Printf("[%v]", err)
}
