package dao

import (
	"fmt"
	"testing"
)

func TestInit(t *testing.T) {
	Init()
}

func TestFind(t *testing.T) {
	Init()
	var tv TableVideo
	result := Db.First(&tv)
	fmt.Println(result.RowsAffected)
	fmt.Println(tv.ID)
	fmt.Println(tv.AuthorId)
	fmt.Println(tv.CoverUrl)
	fmt.Println(tv.PlayUrl)
	fmt.Println(tv.PublishTime)
}

func TestGetVideosByAuthorId(t *testing.T) {
	data, err := GetVideosByAuthorId(2)
	if err != nil {
		print(err)
	}
	for _, video := range data {
		fmt.Println(video)
	}
}
