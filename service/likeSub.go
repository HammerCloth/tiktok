package service

type LikeSub struct {
}

func (like *LikeSub) GetVideo(videoId int64) (Video, error) {
	return Video{
		Id: 1,
		Author: User{
			Id:            1,
			Name:          "lzz",
			FollowCount:   12,
			FollowerCount: 13,
			IsFollow:      true,
		},
		PlayUrl:       "www.baidu.com",
		CoverUrl:      "www.baidu.com",
		FavoriteCount: 2,
		CommentCount:  3,
		IsFavorite:    true,
	}, nil
}
