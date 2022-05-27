package dao

import (
	"fmt"
	"testing"
)

func TestGetTableUserList(t *testing.T) {
	list, err := GetTableUserList()
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestGetTableUserByUsername(t *testing.T) {
	list, err := GetTableUserByUsername("test")
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestGetTableUserById(t *testing.T) {
	list, err := GetTableUserById(int64(4))
	fmt.Printf("%v", list)
	fmt.Printf("%v", err)
}

func TestInsertTableUser(t *testing.T) {
	tu := &TableUser{
		Id:       5,
		Name:     "a",
		Password: "111111",
	}
	list := InsertTableUser(tu)
	fmt.Printf("%v", list)
}
