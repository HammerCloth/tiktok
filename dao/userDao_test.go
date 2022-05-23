package dao

import (
	"fmt"
	"testing"
)

func TestGetTableUserList(t *testing.T) {
	list, err := NewUserDaoInstance().GetTableUserList()
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestGetTableUserByUsername(t *testing.T) {
	list, err := NewUserDaoInstance().GetTableUserByUsername("test")
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestGetTableUserById(t *testing.T) {
	list, err := NewUserDaoInstance().GetTableUserById(int64(4))
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestInsertTableUser(t *testing.T) {
	tu := &TableUser{
		Id:       5,
		Name:     "a",
		Password: "111111",
	}
	list := NewUserDaoInstance().InsertTableUser(tu)
	fmt.Printf("%v", list)
}
