package service

import (
	"fmt"
	"testing"
)

func TestCountComment(t *testing.T) {
	impl := CommentServiceImpl{}
	count, err := impl.CountFromVideoId(1)
	if err != nil {
		fmt.Println("count err:", err)
	}
	fmt.Println(count)
}
