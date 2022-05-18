package service

import (
	"fmt"
	"testing"
	"time"
)

func getVideoService() VideoService {
	var videoService VideoServiceImpl
	videoService.UserService = VideoSub{}
	videoService.LikeService = VideoSub{}
	videoService.CommentService = VideoSub{}
	return &videoService
}

func TestList(t *testing.T) {
	videoService := getVideoService()
	list, err := videoService.List(1)
	if err != nil {
		return
	}
	for _, video := range list {
		fmt.Println(video)
	}

}

func TestGetVideo(t *testing.T) {
	videoService := getVideoService()
	video, err := videoService.GetVideo(1, 1)
	if err != nil {
		return
	}
	fmt.Println(video)
}

func TestFeed(t *testing.T) {
	videoService := getVideoService()
	feed, t2, err := videoService.Feed(time.Now(), 1)
	if err != nil {
		return
	}
	for _, video := range feed {
		fmt.Println(video)
	}
	fmt.Println(t2)
}
