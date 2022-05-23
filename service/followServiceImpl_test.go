package service

import (
	"fmt"
	"testing"
)

func TestIsFollow(t *testing.T) {
	isFollow, err := NewFSIInstance().IsFollowing(1, 2)
	if nil != err {
		t.Errorf("IsFollow() error = %v", err)
	}
	if false == isFollow {
		fmt.Println("不存在该关系")
	}
	fmt.Println("存在该关系")
}
