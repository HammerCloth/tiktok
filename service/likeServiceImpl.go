package service

import (
	"TikTok/config"
	"TikTok/dao"
	"errors"
	"log"
)

type LikeServiceImpl struct {
	VideoService
}

//根据userid,videoid查询点赞信息
func (like *LikeServiceImpl) IsFavourit(videoId int64, userId int64) (bool, error) {
	//查询该点赞信息在数据库中是否存在
	likedata, err := dao.NewLikeDaoInstance().GetLikeInfo(userId, videoId)
	//如果有问题，说明操作数据库失败,返回默认未点赞,输出错误信息err:"get likeInfo failed"
	if err != nil {
		log.Printf("方法GetLikeInfo(userId, videoId) 失败：%v", err)
		return false, err
	} else { //查询数据为0或者查询到数据，根据Cancel值判断是否点赞；
		log.Printf("方法GetLikeInfo(userId, videoId) 成功")
		if likedata == (dao.Like{}) { //查询数据为0
			return false, nil
		} else {
			if likedata.Cancel == config.Islike {
				return true, nil
			} else { //查询cancel为Unlike
				return false, nil
			}
		}

	}
}

//根据videoid获取点赞数量
func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	//查询videoid对应点赞数量
	count, err := dao.NewLikeDaoInstance().GetLikeCount(videoId)
	//如果有问题，说明操作数据库失败,返回默认0,返回错误信息err:"An unknown exception occurred in the query"
	if err != nil {
		log.Printf("方法GetLikeCount(videoId) 失败：%v", err)
		return count, err
	}
	log.Printf("方法GetLikeCount(videoId) 成功")
	return count, err
}

func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, action_type int32) error {
	//如果查询没有数据，用来生成该条点赞信息，存储在likedata中
	var likedata dao.Like
	//先查询是否有这条数据
	likeInfo, err := dao.NewLikeDaoInstance().GetLikeInfo(userId, videoId)
	//点赞行为
	if action_type == config.Likeaction {
		//如果有问题，说明查询数据库失败，返回错误信息err:"get likeInfo failed"
		if err != nil {
			return err
		} else {
			if likeInfo == (dao.Like{}) { //没查到这条数据，则新建这条数据；
				likedata.User_id = userId       //插入userid
				likedata.Video_id = videoId     //插入videoid
				likedata.Cancel = config.Islike //插入点赞cancel=0
				return dao.NewLikeDaoInstance().InsertLike(likedata)
			} else { //查到这条数据,更新即可;
				return dao.NewLikeDaoInstance().UpdateLike(userId, videoId, config.Islike)
			}
		}
	} else { //取消赞行为，只有当前状态是点赞状态才会发起取消赞行为，所以如果查询到，必然是cancel==0(点赞)
		//如果有问题，说明查询数据库失败，返回错误信息err:"get likeInfo failed"
		if err != nil {
			return err
		} else {
			if likeInfo == (dao.Like{}) { //只有当前是点赞状态才能取消点赞这个行为
				// 所以如果查询不到数据则返回错误，err:"can't find data,this action invalid"，就不该有取消赞这个行为
				return errors.New("can't find data,this action invalid")
			} else {
				//如果查询到数据，则更新为取消赞状态
				return dao.NewLikeDaoInstance().UpdateLike(userId, videoId, config.Unlike)
			}
		}
	}
	return nil
}

func (like *LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	//1.先查询点赞列表信息
	likeList, err := dao.NewLikeDaoInstance().GetLikeList(userId)
	//如果有问题，说明查询数据库失败，返回空和错误err:"get likeList failed"
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	//提前定义好切片长度,生成集合
	favorite_videolist := make([]Video, 0, len(likeList))
	//如果查询成功，无论是否有数据，遍历likelist,获得其中的video_id；
	//测试结构体，协同开发
	//likesub := new(LikeSub)
	for _, likedata := range likeList {
		//测试函数，协同开发
		//video, err1 := likesub.GetVideo(likedata.Video_id,userId)
		//调用video接口，Getvideo：根据videoid，当前用户id，返回video对象
		video, err1 := like.GetVideo(likedata.Video_id, curId)
		if err1 != nil { //如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且跳过
			log.Println(errors.New("can't find this favourite video"))
			continue
		} //将每个video对象添加到集合中去
		favorite_videolist = append(favorite_videolist, video)
	}
	return favorite_videolist, nil
}
