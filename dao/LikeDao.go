package dao

import (
	"TikTok/config"
	"errors"
	"log"
)

//like表的结构。
type Like struct {
	Id       int64 //自增主键
	User_id  int64 //点赞用户id
	Video_id int64 //视频id
	Cancel   int8  //是否点赞，0为点赞，1为取消赞
}

// TableName 修改表名映射
func (Like) TableName() string {
	return "likes"
}

//1.根据videoid获取点赞userid
func GetLikeUserIdList(videoId int64) ([]int64, error) {
	var likeUserIdList []int64 //存所有该视频点赞用户id；
	//查询likes表对应视频id点赞用户，返回查询结果
	err := Db.Model(Like{}).Where(map[string]interface{}{"video_id": videoId, "cancel": config.Islike}).
		Pluck("user_id", &likeUserIdList).Error
	//查询过程出现错误，返回默认值0，并输出错误信息
	if err != nil {
		log.Println(err.Error())
		return nil, errors.New("An unknown exception occurred in the query")
	} else {
		//没查询到或者查询到结果，返回数量以及无报错
		return likeUserIdList, nil
	}
}

//2.根据userid，videoid,action_type点赞或者取消赞
func UpdateLike(userId int64, videoId int64, action_type int32) error {
	//更新当前用户观看视频的点赞状态“cancel”，返回错误结果
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		Update("cancel", action_type).Error
	//如果出现错误，返回更新数据库失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("update data fail")
	}
	//更新操作成功
	return nil
}

//3、插入点赞数据
func InsertLike(likedata Like) error {
	//创建点赞数据，默认为点赞，cancel为0，返回错误结果
	err := Db.Model(Like{}).Create(&likedata).Error
	//如果有错误结果，返回插入失败
	if err != nil {
		log.Println(err.Error())
		return errors.New("insert data fail")
	}
	return nil
}

//4.根据userid,videoid查询点赞信息
func GetLikeInfo(userId int64, videoId int64) (Like, error) {
	//创建一条空like结构体，用来存储查询到的信息
	var likeInfo Like
	//根据userid,videoid查询是否有该条信息，如果有，存储在likeInfo,返回查询结果
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "video_id": videoId}).
		First(&likeInfo).Error
	if err != nil {
		//查询数据为0，打印"can't find data"，返回空结构体，这时候就应该要考虑是否插入这条数据了
		if "record not found" == err.Error() {
			log.Println("can't find data")
			return Like{}, nil
		} else {
			//如果查询数据库失败，返回获取likeInfo信息失败
			log.Println(err.Error())
			return likeInfo, errors.New("get likeInfo failed")
		}
	}
	return likeInfo, nil
}

////5.根据userid查询所属点赞全部列表信息
//func GetLikeList(userId int64) ([]Like, error) {
//	//创建likeList切片，用来存储查询到的当前用户点赞列表信息
//	var likeList []Like
//	//根据userid查询所有点赞视频信息，如果有，存储在likeList切片中,返回查询结果
//	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "cancel": config.Islike}).
//		Find(&likeList).Error
//	if err != nil {
//		//查询数据为0，返回空likeList切片，以及返回无错误
//		if "record not found" == err.Error() {
//			log.Println("there are no likes")
//			return likeList, nil
//		} else {
//			//如果查询数据库失败，返回获取likelist失败
//			log.Println(err.Error())
//			return likeList, errors.New("get likeList failed")
//		}
//	}
//	return likeList, nil
//}

//6.根据userid查询所属点赞全部videoid
func GetLikeVideoIdList(userId int64) ([]int64, error) {
	var likeVideoIdList []int64
	err := Db.Model(Like{}).Where(map[string]interface{}{"user_id": userId, "cancel": config.Islike}).
		Pluck("video_id", &likeVideoIdList).Error
	if err != nil {
		//查询数据为0，返回空likeVideoIdList切片，以及返回无错误
		if "record not found" == err.Error() {
			log.Println("there are no likesvideoid")
			return likeVideoIdList, nil
		} else {
			//如果查询数据库失败，返回获取likeVideoIdList失败
			log.Println(err.Error())
			return likeVideoIdList, errors.New("get likeVideoIdList failed")
		}
	}
	return likeVideoIdList, nil
}
