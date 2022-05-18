package service

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	var videoService VideoServiceImpl
	list, err := videoService.List(2)
	if err != nil {
		return
	}
	for _, video := range list {
		fmt.Println(video)
	}

}

func TestGetVideo(t *testing.T) {
	var videoService VideoServiceImpl
	videoService.UserService = VideoSub{}
	video, err := videoService.GetVideo(1, 1)
	if err != nil {
		return
	}
	fmt.Println(video)

}
