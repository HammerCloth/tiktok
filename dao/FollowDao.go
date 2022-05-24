package dao

import (
	"log"
	"sync"
)

// Follow 用户关系结构，对应用户关系表。
type Follow struct {
	Id         int64
	UserId     int64
	FollowerId int64
	Cancel     int8
}

// TableName 设置Follow结构体对应数据库表名。
func (Follow) TableName() string {
	return "follows"
}

// FollowDao 把dao层看成整体，把dao的curd封装在一个结构体中。
type FollowDao struct {
}

var (
	followDao  *FollowDao //操作该dao层crud的结构体变量。
	followOnce sync.Once  //单例限定，去限定申请一个followDao结构体变量。
)

// NewFollowDaoInstance 生成并返回followDao的单例对象。
func NewFollowDaoInstance() *FollowDao {
	followOnce.Do(
		func() {
			followDao = &FollowDao{}
		})
	return followDao
}

/*
下面为FollowDao的成员方法，即crud逻辑。
*/

// FindRelation 给定当前用户和目标用户id，查询follow表中相应的记录。
func (*FollowDao) FindRelation(userId int64, targetId int64) (*Follow, error) {
	// follow变量用于后续存储数据库查出来的用户关系。
	follow := Follow{}
	//当查询出现错误时，日志打印err msg，并return err.
	if err := Db.
		Where("user_id = ?", targetId).
		Where("follower_id = ?", userId).
		Where("cancel = ?", 0).
		Take(&follow).Error; nil != err {
		// 当没查到数据时，gorm也会报错。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}
	//正常情况，返回取到的值和空err.
	return &follow, nil
}

// GetFollowerCnt 给定当前用户id，查询follow表中该用户的粉丝数。
func (*FollowDao) GetFollowerCnt(userId int64) (int64, error) {
	// 用于存储当前用户粉丝数的变量
	var cnt int64
	// 当查询出现错误的情况，日志打印err msg，并返回err.
	if err := Db.
		Model(Follow{}).
		Where("user_id = ?", userId).
		Where("cancel = ?", 0).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	// 正常情况，返回取到的粉丝数。
	return cnt, nil
}

// GetFollowingCnt 给定当前用户id，查询follow表中该用户关注了多少人。
func (*FollowDao) GetFollowingCnt(userId int64) (int64, error) {
	// 用于存储当前用户关注了多少人。
	var cnt int64
	// 查询出错，日志打印err msg，并return err
	if err := Db.Model(Follow{}).
		Where("follower_id = ?", userId).
		Where("cancel = ?", 0).
		Count(&cnt).Error; nil != err {
		log.Println(err.Error())
		return 0, err
	}
	// 查询成功，返回人数。
	return cnt, nil
}

// InsertFollowRelation 给定用户和目标对象id，插入其关注关系。
func (*FollowDao) InsertFollowRelation(userId int64, targetId int64) (bool, error) {
	// 生成需要插入的关系结构体。
	follow := Follow{
		UserId:     userId,
		FollowerId: targetId,
		Cancel:     0,
	}
	// 插入失败，返回err.
	if err := Db.Select("UserId", "FollowerId", "Cancel").Create(&follow).Error; nil != err {
		log.Println(err.Error())
		return false, err
	}
	// 插入成功
	return true, nil
}

// FindEverFollowing 给定当前用户和目标用户id，查看曾经是否有关注关系。
func (*FollowDao) FindEverFollowing(userId int64, targetId int64) (*Follow, error) {
	// 用于存储查出来的关注关系。
	follow := Follow{}
	//当查询出现错误时，日志打印err msg，并return err.
	if err := Db.
		Where("user_id = ?", userId).
		Where("follower_id = ?", targetId).
		Where("cancel = ? or cancel = ?", 0, 1).
		Take(&follow).Error; nil != err {
		// 当没查到记录报错时，不当做错误处理。
		if "record not found" == err.Error() {
			return nil, nil
		}
		log.Println(err.Error())
		return nil, err
	}
	//正常情况，返回取到的关系和空err.
	return &follow, nil
}

// UpdateFollowRelation 给定用户和目标用户的id，更新他们的关系为取消关注或再次关注。
func (*FollowDao) UpdateFollowRelation(userId int64, targetId int64, cancel int8) (bool, error) {
	// 更新失败，返回错误。
	if err := Db.Model(Follow{}).
		Where("user_id = ?", userId).
		Where("follower_id = ?", targetId).
		Update("cancel", cancel).Error; nil != err {
		// 更新失败，打印错误日志。
		log.Println(err.Error())
		return false, err
	}
	// 更新成功。
	return true, nil
}

// GetFollowingIds 给定用户id，查询他关注了哪些人的id。
func (*FollowDao) GetFollowingIds(userId int64) ([]int64, error) {
	var ids []int64
	if err := Db.
		Model(Follow{}).
		Where("follower_id = ?", userId).
		Pluck("user_id", &ids).Error; nil != err {
		// 没有关注任何人，但是不能算错。
		if "record not found" == err.Error() {
			return nil, nil
		}
		// 查询出错。
		log.Println(err.Error())
		return nil, err
	}
	// 查询成功。
	return ids, nil
}

// GetFollowersIds 给定用户id，查询他关注了哪些人的id。
func (*FollowDao) GetFollowersIds(userId int64) ([]int64, error) {
	var ids []int64
	if err := Db.
		Model(Follow{}).
		Where("user_id = ?", userId).
		Where("cancel = ?", 0).
		Pluck("follower_id", &ids).Error; nil != err {
		// 没有粉丝，但是不能算错。
		if "record not found" == err.Error() {
			return nil, nil
		}
		// 查询出错。
		log.Println(err.Error())
		return nil, err
	}
	// 查询成功。
	return ids, nil
}
