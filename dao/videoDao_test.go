package dao

import (
	"fmt"
	"testing"
	"time"
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

func TestGetVideoByVideoId(t *testing.T) {
	data, err := GetVideoByVideoId(1)
	if err != nil {
		print(err)
	}
	fmt.Println(data)

}

func TestGetVideosByLastTime(t *testing.T) {
	data, err := GetVideosByLastTime(time.Now())
	if err != nil {
		return
	}
	for _, video := range data {
		fmt.Println(video)
	}
}
