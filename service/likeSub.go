package service

type LikeSub struct {
}

func (like *LikeSub) GetVideo(videoId int64, userId int64) (Video, error) {
	if videoId%2 == 1 {
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
	} else {
		return Video{
			Id: 2,
			Author: User{
				Id:            2,
				Name:          "lzz2",
				FollowCount:   11,
				FollowerCount: 11,
				IsFollow:      false,
			},
			PlayUrl:       "www.baidu11.com",
			CoverUrl:      "www.baidu11.com",
			FavoriteCount: 1,
			CommentCount:  3,
			IsFavorite:    false,
		}, nil
	}

}
