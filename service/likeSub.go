package service

import "TikTok/dao"

type LikeSub struct {
}

func (like *LikeSub) GetVideo(videoId int64, userId int64) (Video, error) {
	if videoId%2 == 1 {
		return Video{
			TableVideo: dao.TableVideo{
				Id:       1,
				AuthorId: 1,
				PlayUrl:  "www.baidu.com",
				CoverUrl: "www.baidu.com",
			},
			Author: User{
				Id:            1,
				Name:          "lzz",
				FollowCount:   12,
				FollowerCount: 13,
				IsFollow:      true,
			},
			FavoriteCount: 2,
			CommentCount:  3,
			IsFavorite:    true,
		}, nil
	} else {
		return Video{
			TableVideo: dao.TableVideo{
				Id:       2,
				AuthorId: 3,
				PlayUrl:  "www.baidu.com",
				CoverUrl: "www.baidu11.com",
			},
			Author: User{
				Id:            2,
				Name:          "lzz",
				FollowCount:   12,
				FollowerCount: 13,
				IsFollow:      true,
			},
			FavoriteCount: 2,
			CommentCount:  3,
			IsFavorite:    true,
		}, nil
	}

}
