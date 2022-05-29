package service

import (
	"TikTok/dao"
	"TikTok/middleware"
	"log"
	"strconv"
	"strings"
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
	// 先查Redis里面是否有此关系。
	if flag, err := middleware.RdbFollowingPart.SIsMember(middleware.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
		return true, err
	}
	// SQL 查询。
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
	// 查Redis中是否已经存在。
	if cnt, err := middleware.RdbFollowers.SCard(middleware.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		return cnt, err
	}
	// SQL中查询。
	cnt, err := dao.NewFollowDaoInstance().GetFollowerCnt(userId)

	if nil != err {
		return 0, err
	}
	return cnt, err
}

// GetFollowingCnt 给定当前用户id，查询其关注者数量。
func (*FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	// 查看Redis中是否有关注数。
	if cnt, err := middleware.RdbFollowing.SCard(middleware.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		return cnt, err
	}
	// 用SQL查询。
	cnt, err := dao.NewFollowDaoInstance().GetFollowingCnt(userId)

	if nil != err {
		return 0, err
	}

	return cnt, err
}

// AddFollowRelation 给定当前用户和目标对象id，添加他们之间的关注关系。
func (*FollowServiceImp) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	middleware.RmqFollowAdd.Publish(sb.String())
	// 记录日志
	log.Println("消息打入成功。")
	// 更新redis信息。
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		3-Redis是否存在user_targetId和user_userId.
		4-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := middleware.RdbFollowers.SCard(middleware.Ctx, targetIdStr).Result(); 0 != cnt {
		middleware.RdbFollowers.SAdd(middleware.Ctx, targetIdStr, userId)
	}
	// step2
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbFollowing.SCard(middleware.Ctx, followingIdStr).Result(); 0 != cnt {
		middleware.RdbFollowing.SAdd(middleware.Ctx, followingIdStr, targetId)
	}
	// step3
	userTargetIdStr := strconv.Itoa(int(targetId))
	userUserIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbUser.Exists(middleware.Ctx, userTargetIdStr).Result(); cnt > 0 {
		param, _ := middleware.RdbUser.HGet(middleware.Ctx, userTargetIdStr, "follower_count").Result()
		followerCount, _ := strconv.Atoi(param)
		middleware.RdbUser.HSet(middleware.Ctx, userTargetIdStr, "follower_count", followerCount+1)
	}
	if cnt, _ := middleware.RdbUser.Exists(middleware.Ctx, userUserIdStr).Result(); cnt > 0 {
		param, _ := middleware.RdbUser.HGet(middleware.Ctx, userUserIdStr, "follow_count").Result()
		followCount, _ := strconv.Atoi(param)
		middleware.RdbUser.HSet(middleware.Ctx, userUserIdStr, "follow_count", followCount+1)
	}
	// step4
	followingPartUserIdStr := strconv.Itoa(int(userId))
	middleware.RdbFollowingPart.SAdd(middleware.Ctx, followingPartUserIdStr, targetId)
	return true, nil
	/*followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(targetId, userId)
	// 寻找SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(targetId, userId, 0)
		// update 出错。
		if nil != err {
			return false, err
		}
		// update 成功。
		return true, nil
	}
	// 曾经没有关注过，需要插入一条关注关系。
	_, err = followDao.InsertFollowRelation(targetId, userId)
	if nil != err {
		// insert 出错
		return false, err
	}
	// insert 成功。
	return true, nil*/
}

// DeleteFollowRelation 给定当前用户和目标用户id，删除其关注关系。
func (*FollowServiceImp) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	middleware.RmqFollowDel.Publish(sb.String())
	// 记录日志
	log.Println("消息打入成功。")
	// 更新redis信息。
	return updateRedisWithDel(userId, targetId)
	/*followDao := dao.NewFollowDaoInstance()
	follow, err := followDao.FindEverFollowing(targetId, userId)
	// 寻找 SQL 出错。
	if nil != err {
		return false, err
	}
	// 曾经关注过，只需要update一下cancel即可。
	if nil != follow {
		_, err := followDao.UpdateFollowRelation(targetId, userId, 1)
		// update 出错。
		if nil != err {
			return false, err
		}
		// update 成功。
		return true, nil
	}
	// 没有关注关系
	return false, nil*/
}

// 当取关时，更新redis里的信息
func updateRedisWithDel(userId int64, targetId int64) (bool, error) {
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		3-Redis是否存在user_targetId和user_userId.
		4-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := middleware.RdbFollowers.SCard(middleware.Ctx, targetIdStr).Result(); 0 != cnt {
		middleware.RdbFollowers.SRem(middleware.Ctx, targetIdStr, userId)
	}
	// step2
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbFollowing.SCard(middleware.Ctx, followingIdStr).Result(); 0 != cnt {
		middleware.RdbFollowing.SRem(middleware.Ctx, followingIdStr, targetId)
	}
	// step3
	userTargetIdStr := strconv.Itoa(int(targetId))
	userUserIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbUser.Exists(middleware.Ctx, userTargetIdStr).Result(); cnt > 0 {
		param, _ := middleware.RdbUser.HGet(middleware.Ctx, userTargetIdStr, "follower_count").Result()
		followerCount, _ := strconv.Atoi(param)
		middleware.RdbUser.HSet(middleware.Ctx, userTargetIdStr, "follower_count", followerCount-1)
	}
	if cnt, _ := middleware.RdbUser.Exists(middleware.Ctx, userUserIdStr, "name").Result(); cnt > 0 {
		param, _ := middleware.RdbUser.HGet(middleware.Ctx, userUserIdStr, "follow_count").Result()
		followCount, _ := strconv.Atoi(param)
		middleware.RdbUser.HSet(middleware.Ctx, userUserIdStr, "follow_count", followCount-1)
	}
	// step4
	followingPartUserIdStr := strconv.Itoa(int(userId))
	middleware.RdbFollowingPart.SRem(middleware.Ctx, followingPartUserIdStr, targetId)
	return true, nil
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
	// 先查Redis，看是否有全部关注信息。
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbFollowers.SCard(middleware.Ctx, followingIdStr).Result(); 0 == cnt {
		users, _ := getFollowing(userId)

		go setRedisFollowing(userId, users)

		return users, nil
	}
	// Redis中有。
	UserIdStr := strconv.Itoa(int(userId))
	userIds, _ := middleware.RdbFollowing.SMembers(middleware.Ctx, UserIdStr).Result()
	len := len(userIds)
	users := make([]User, len)
	wg := sync.WaitGroup{}
	for i := 0; i < len; i++ {
		wg.Add(1)
		go func(k int) {
			defer wg.Done()
			user, _ := middleware.RdbUser.HGetAll(middleware.Ctx, userIds[k]).Result()
			uid, _ := strconv.Atoi(userIds[k])
			users[k].Id = int64(uid)
			users[k].Name = user["name"]
			cnt, _ := strconv.Atoi(user["follow_count"])
			users[k].FollowCount = int64(cnt)
			cnt, _ = strconv.Atoi(user["follower_count"])
			users[k].FollowerCount = int64(cnt)
			// 必然是关注情况。
			users[k].IsFollow = true
			//}
		}(i)
	}
	wg.Wait()
	log.Println("从Redis中查询到所有关注者。")
	return users, nil
}

// 设置Redis关于所有关注的信息。
func setRedisFollowing(userId int64, users []User) {
	/*
		1-设置following_userId的所有关注id。
		2-设置user_userId基本信息。
		3-设置following_part_id关注信息。
	*/
	for _, user := range users {

		followingIdStr := strconv.Itoa(int(userId))
		middleware.RdbFollowing.SAdd(middleware.Ctx, followingIdStr, user.Id)
		userUserIdStr := strconv.Itoa(int(user.Id))
		middleware.RdbUser.HSet(middleware.Ctx, userUserIdStr, map[string]interface{}{
			"name":           user.Name,
			"follow_count":   user.FollowCount,
			"follower_count": user.FollowerCount,
		})
		middleware.RdbFollowingPart.SAdd(middleware.Ctx, followingIdStr, user.Id)
	}
}

// 从数据库查所有关注用户信息。
func getFollowing(userId int64) ([]User, error) {
	users := make([]User, 1)
	// 查询出错。
	if err := dao.Db.Raw("select id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_count,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_count,"+
		"\n 'true' is_follow\nfrom\n("+
		"\nselect f1.follower_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0\n\tunion all"+
		"\nselect f1.follower_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.user_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0\n) T"+
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
	// 先查Redis，看是否有全部粉丝信息。
	followersIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbFollowers.SCard(middleware.Ctx, followersIdStr).Result(); 0 == cnt {
		users, _ := getFollowers(userId)

		go setRedisFollowers(userId, users)

		return users, nil
	}
	// Redis中有。
	userIds, _ := middleware.RdbFollowers.SMembers(middleware.Ctx, followersIdStr).Result()
	len := len(userIds)
	users := make([]User, len)
	wg := sync.WaitGroup{}
	for i := 0; i < len; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			user, _ := middleware.RdbUser.HGetAll(middleware.Ctx, userIds[i]).Result()
			uid, _ := strconv.Atoi(userIds[i])
			users[i].Id = int64(uid)
			users[i].Name = user["name"]
			cnt, _ := strconv.Atoi(user["follow_count"])
			users[i].FollowCount = int64(cnt)
			cnt, _ = strconv.Atoi(user["follower_count"])
			users[i].FollowerCount = int64(cnt)
			// 从following_part_#{id}中看是否有关注关系。
			isFollow, _ := middleware.RdbFollowingPart.SIsMember(middleware.Ctx, followersIdStr, userIds[i]).Result()
			users[i].IsFollow = isFollow
		}(i)
	}
	wg.Wait()
	// log.Println("从Redis中查询到所有粉丝们。")
	return users, nil
}

// 重数据库查所有粉丝信息。
func getFollowers(userId int64) ([]User, error) {
	users := make([]User, 1)

	if err := dao.Db.Raw("select T.id,T.name,T.follow_cnt follow_count,T.follower_cnt follower_count,if(f.cancel is null,'false','true') is_follow"+
		"\nfrom follows f right join"+
		"\n(select fid,id,`name`,"+
		"\ncount(if(tag = 'follower' and cancel is not null,1,null)) follower_cnt,"+
		"\ncount(if(tag = 'follow' and cancel is not null,1,null)) follow_cnt"+
		"\nfrom("+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follower' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.user_id and f2.cancel = 0"+
		"\nunion all"+
		"\nselect f1.user_id fid,u.id,`name`,f2.cancel,'follow' tag"+
		"\nfrom follows f1 join users u on f1.follower_id = u.id and f1.cancel = 0"+
		"\nleft join follows f2 on u.id = f2.follower_id and f2.cancel = 0"+
		"\n) T group by fid,id,`name`"+
		"\n) T on f.user_id = T.id and f.follower_id = T.fid and f.cancel = 0 where fid = ?", userId).
		Scan(&users).Error; nil != err {
		// 查询出错。
		return nil, err
	}
	// 查询成功。
	return users, nil
}

// 设置Redis关于所有粉丝的信息
func setRedisFollowers(userId int64, users []User) {
	/*
		1-设置followers_userId的所有粉丝id。
		2-设置user_userId基本信息。
		3-设置following_part_id关注信息。
	*/
	for _, user := range users {
		followersIdStr := strconv.Itoa(int(userId))
		middleware.RdbFollowers.SAdd(middleware.Ctx, followersIdStr, user.Id)
		userUserIdStr := strconv.Itoa(int(user.Id))
		middleware.RdbUser.HSet(middleware.Ctx, userUserIdStr, map[string]interface{}{
			"name":           user.Name,
			"follow_count":   user.FollowCount,
			"follower_count": user.FollowerCount,
		})
		middleware.RdbFollowingPart.SAdd(middleware.Ctx, userUserIdStr, userId)

		if user.IsFollow {
			middleware.RdbFollowingPart.SAdd(middleware.Ctx, followersIdStr, user.Id)
		}
	}
}
