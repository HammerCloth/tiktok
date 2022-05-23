package service

import (
	"TikTok/dao"
	"sync"
)

// FollowServiceImp 该结构体继承FollowService接口。
type FollowServiceImp struct {
	UserService
}

var (
	followServiceImp  *FollowServiceImp //controller层通过该实例变量调用service的所有业务方法。
	followServiceOnce sync.Once         //限定该service对象为单例，节约内存。
)

// NewFSIInstance 生成并返回FollowServiceImp结构体单例变量。
func NewFSIInstance() *FollowServiceImp {
	followServiceOnce.Do(
		func() {
			followServiceImp = &FollowServiceImp{
				UserService: &UserServiceImpl{
					// 存在我调userService中，userService又要调我。
					FollowService: &FollowServiceImp{},
				},
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

// GetFollowerCnt 给定当前用户id，查询其粉丝数量。
func (*FollowServiceImp) GetFollowerCnt(userId int64) (int64, error) {
	cnt, err := dao.NewFollowDaoInstance().GetFollowerCnt(userId)

	if nil != err {
		return 0, err
	}
	return cnt, err
}

// GetFollowingCnt 给定当前用户id，查询其关注者数量。
func (*FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	cnt, err := dao.NewFollowDaoInstance().GetFollowingCnt(userId)

	if nil != err {
		return 0, err
	}

	return cnt, err
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
		_, err := followDao.UpdateFollowRelation(userId, targetId, 0)
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
		_, err := followDao.UpdateFollowRelation(userId, targetId, 1)
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
/*func (f *FollowServiceImp) GetFollowing(userId int64) ([]User, error) {
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
}*/
// GetFollowing 根据当前用户id来查询他的关注者列表。
func (f *FollowServiceImp) GetFollowing(userId int64) ([]User, error) {
	users := make([]User, 1)
	// 查询出错。
	if err := dao.Db.Raw("select id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_count,"+
		"\n'true' isFollow\nfrom\n("+
		"\n\tselect f1.follower_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\n\tfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\n\tleft join follows f2 on u.id = f2.user_id and f2.cancel = 0\n\tunion all"+
		"\n\tselect f1.follower_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\n\tfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\n\tleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0\n) T"+
		"\nwhere fid = ? group by fid,id,`name`", userId).Scan(&users).Error; nil != err {
		return nil, err
	}
	// 返回关注对象列表。
	return users, nil
}

// GetFollowers 根据当前用户id来查询他的粉丝列表。

/*func (f *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
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
}*/

// GetFollowers 根据当前用户id来查询他的粉丝列表。
func (f *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	users := make([]User, 1)

	if err := dao.Db.Raw("select T.id,T.name,T.follow_cnt follow_count,T.follower_cnt follower_count,if(f.cancel is null,'false','true') is_follow"+
		"\nfrom follows f right join"+
		"\n(\n\tselect fid,id,`name`,"+
		"\n\tcount(if(tag = 'follower' and cancel is not null,1,null)) follower_cnt,"+
		"\n\tcount(if(tag = 'follow' and cancel is not null,1,null)) follow_cnt"+
		"\n\tfrom\n\t\t("+
		"\n\t\tselect f1.user_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\n\t\tfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\n\t\tleft join follows f2 on u.id = f2.user_id and f2.cancel = 0"+
		"\n\t\tunion all"+
		"\n\t\tselect f1.user_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\n\t\tfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\n\t\tleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0"+
		"\n\t\t) T\n\t\tgroup by fid,id,`name`"+
		"\n) T on f.user_id = T.fid and f.follower_id = T.id and f.cancel = 0 where fid = ?", userId).
		Scan(&users).Error; nil != err {
		// 查询出错。
		return nil, err
	}
	// 查询成功。
	return users, nil
}
