package service

import (
	"TikTok/dao"
	"sync"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
	UserService
	FollowService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce *sync.Once        //限定该service对象为单例，节约内存。
)

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowServiceImp {
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{
				UserService: &UserServiceImpl{},
			}
		})
	return followServiceImp
}

// IsFollowing 给定当前用户和目标用户id，判断是否存在关注关系。
func (*FollowServiceImp) IsFollowing(userId int64, targetId int64) (bool, error) {
	relation, err := dao.NewFollowDaoInstance().FindRelation(userId, targetId)

	if nil != err {
		return false, err
	}
	if nil == relation {
		return false, nil
	}
	return true, nil
}

// AddFollowRelation 给定当前用户和目标对象id，添加他们之间的关注关系。
func (*FollowServiceImp) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(userId, targetId, 1)
		// update 出错。
		if nil != err {
			return false, err
		}
		// update 成功。
		return true, nil
	}
	// 曾经没有关注过，需要插入一条关注关系。
	_, err = followDao.InsertFollowRelation(userId, targetId)
	if nil != err {
		// insert 出错
		return false, err
	}
	// insert 成功。
	return true, nil
}

// DeleteFollowRelation 给定当前用户和目标用户id，删除其关注关系。
func (*FollowServiceImp) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(userId, targetId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(userId, targetId, 0)
		// update 出错。
		if nil != err {
			return false, err
		}
		// update 成功。
		return true, nil
	}
	// 没有关注关系
	return false, nil
}

// GetFollowing 根据当前用户id来查询他的关注者列表。
func (f *FollowServiceImp) GetFollowing(userId int64) ([]User, error) {
	// 获取关注对象的id数组。
	ids, err := dao.NewFollowDaoInstance().GetFollowingIds(userId)
	// 查询出错
	if nil != err {
		return nil, err
	}
	// 没得关注者
	if nil == ids {
		return nil, nil
	}
	// 根据每个id来查询用户信息。
	length := len(ids)
	users := make([]User, length)
	for i := 0; i < length; i++ {
		user, err := f.GetUserById(ids[i])
		// 查询失败，继续查其他的。
		if nil != err {
			continue
		}
		// 查询成功，把user加入相应的位置。
		users[i] = user
	}
	// 返回关注对象列表。
	return users, nil
}

// GetFollowers 根据当前用户id来查询他的粉丝列表。
func (f *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	// 获取粉丝的id数组。
	ids, err := dao.NewFollowDaoInstance().GetFollowersIds(userId)
	// 查询出错
	if nil != err {
		return nil, err
	}
	// 没得粉丝
	if nil == ids {
		return nil, nil
	}
	// 根据每个id来查询用户信息。
	length := len(ids)
	users := make([]User, length)
	for i := 0; i < length; i++ {
		user, err := f.GetUserById(ids[i])
		// 查询失败，继续查其他的。
		if nil != err {
			continue
		}
		// 查询成功，把user加入相应的位置。
		users[i] = user
	}
	// 返回粉丝列表。
	return users, nil
}
