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

func (like *LikeSub) GetVideoIdList(userId int64) ([]int64, error) {
	videoList := []int64{51, 52, 53, 54, 55, 56, 57, 58, 59, 60, 61, 62, 63, 64, 65,
		66, 67, 68, 69, 70, 71, 72, 73, 74, 75}
	return videoList, nil
}
