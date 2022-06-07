package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/rabbitmq"
	"TikTok/middleware/redis"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
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
	if flag, err := redis.RdbFollowingPart.SIsMember(redis.Ctx, strconv.Itoa(int(userId)), targetId).Result(); flag {
		// 重现设置过期时间。
		redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
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
	// 存在此关系，将其注入Redis中。
	go addRelationToRedis(int(userId), int(targetId))

	return true, nil
}
func addRelationToRedis(userId int, targetId int) {
	// 第一次存入时，给该key添加一个-1为key，防止脏数据的写入。当然set可以去重，直接加，便于CPU。
	redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(userId)), -1)
	// 将查询到的关注关系注入Redis.
	redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(userId)), targetId)
	// 更新过期时间。
	redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
}

// GetFollowerCnt 给定当前用户id，查询其粉丝数量。
func (*FollowServiceImp) GetFollowerCnt(userId int64) (int64, error) {
	// 查Redis中是否已经存在。
	if cnt, err := redis.RdbFollowers.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}
	// SQL中查询。
	ids, err := dao.NewFollowDaoInstance().GetFollowersIds(userId)
	if nil != err {
		return 0, err
	}
	// 将数据存入Redis.
	// 更新followers 和 followingPart
	go addFollowersToRedis(int(userId), ids)

	return int64(len(ids)), err
}
func addFollowersToRedis(userId int, ids []int64) {
	redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	for i, id := range ids {
		redis.RdbFollowers.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(id)), userId)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(int(id)), -1)
		// 更新部分关注者的时间
		redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(int(id)),
			config.ExpireTime+time.Duration((i%10)<<8))
	}
	// 更新followers的过期时间。
	redis.RdbFollowers.Expire(redis.Ctx, strconv.Itoa(userId), config.ExpireTime)

}

// GetFollowingCnt 给定当前用户id，查询其关注者数量。
func (*FollowServiceImp) GetFollowingCnt(userId int64) (int64, error) {
	// 查看Redis中是否有关注数。
	if cnt, err := redis.RdbFollowing.SCard(redis.Ctx, strconv.Itoa(int(userId))).Result(); cnt > 0 {
		// 更新过期时间。
		redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(int(userId)), config.ExpireTime)
		return cnt - 1, err
	}
	// 用SQL查询。
	ids, err := dao.NewFollowDaoInstance().GetFollowingIds(userId)

	if nil != err {
		return 0, err
	}
	// 更新Redis中的followers和followPart
	go addFollowingToRedis(int(userId), ids)

	return int64(len(ids)), err
}
func addFollowingToRedis(userId int, ids []int64) {
	redis.RdbFollowing.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
	for i, id := range ids {
		redis.RdbFollowing.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(userId), id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, strconv.Itoa(userId), -1)
		// 更新过期时间
		redis.RdbFollowingPart.Expire(redis.Ctx, strconv.Itoa(userId),
			config.ExpireTime+time.Duration((i%10)<<8))
	}
	// 更新following的过期时间
	redis.RdbFollowing.Expire(redis.Ctx, strconv.Itoa(userId), config.ExpireTime)
}

// AddFollowRelation 给定当前用户和目标对象id，添加他们之间的关注关系。
func (*FollowServiceImp) AddFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	rabbitmq.RmqFollowAdd.Publish(sb.String())
	// 记录日志
	log.Println("消息打入成功。")
	// 更新redis信息。
	return updateRedisWithAdd(userId, targetId)
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

// 添加关注时，设置Redis
func updateRedisWithAdd(userId int64, targetId int64) (bool, error) {
	/*
		1-Redis是否存在followers_targetId.
		2-Redis是否存在following_userId.
		3-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := redis.RdbFollowers.SCard(redis.Ctx, targetIdStr).Result(); 0 != cnt {
		redis.RdbFollowers.SAdd(redis.Ctx, targetIdStr, userId)
		redis.RdbFollowers.Expire(redis.Ctx, targetIdStr, config.ExpireTime)
	}
	// step2
	followingUserIdStr := strconv.Itoa(int(userId))
	if cnt, _ := redis.RdbFollowing.SCard(redis.Ctx, followingUserIdStr).Result(); 0 != cnt {
		redis.RdbFollowing.SAdd(redis.Ctx, followingUserIdStr, targetId)
		redis.RdbFollowing.Expire(redis.Ctx, followingUserIdStr, config.ExpireTime)
	}
	// step3
	followingPartUserIdStr := followingUserIdStr
	redis.RdbFollowingPart.SAdd(redis.Ctx, followingPartUserIdStr, targetId)
	// 可能是第一次给改用户加followingPart的关注者，需要加上-1防止脏读。
	redis.RdbFollowingPart.SAdd(redis.Ctx, followingPartUserIdStr, -1)
	redis.RdbFollowingPart.Expire(redis.Ctx, followingPartUserIdStr, config.ExpireTime)
	return true, nil
}

// DeleteFollowRelation 给定当前用户和目标用户id，删除其关注关系。
func (*FollowServiceImp) DeleteFollowRelation(userId int64, targetId int64) (bool, error) {
	// 加信息打入消息队列。
	sb := strings.Builder{}
	sb.WriteString(strconv.Itoa(int(userId)))
	sb.WriteString(" ")
	sb.WriteString(strconv.Itoa(int(targetId)))
	rabbitmq.RmqFollowDel.Publish(sb.String())
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
		2-Redis是否存在following_part_userId.
	*/
	// step1
	targetIdStr := strconv.Itoa(int(targetId))
	if cnt, _ := redis.RdbFollowers.SCard(redis.Ctx, targetIdStr).Result(); 0 != cnt {
		redis.RdbFollowers.SRem(redis.Ctx, targetIdStr, userId)
		redis.RdbFollowers.Expire(redis.Ctx, targetIdStr, config.ExpireTime)
	}
	// step2
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := redis.RdbFollowing.SCard(redis.Ctx, followingIdStr).Result(); 0 != cnt {
		redis.RdbFollowing.SRem(redis.Ctx, followingIdStr, targetId)
		redis.RdbFollowing.Expire(redis.Ctx, followingIdStr, config.ExpireTime)
	}
	// step3
	followingPartUserIdStr := followingIdStr
	if cnt, _ := redis.RdbFollowingPart.Exists(redis.Ctx, followingPartUserIdStr).Result(); 0 != cnt {
		redis.RdbFollowingPart.SRem(redis.Ctx, followingPartUserIdStr, targetId)
		redis.RdbFollowingPart.Expire(redis.Ctx, followingPartUserIdStr, config.ExpireTime)
	}
	return true, nil
}

// GetFollowing 根据当前用户id来查询他的关注者列表。
func (f *FollowServiceImp) getFollowing(userId int64) ([]User, error) {
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
	len := len(ids)
	if len > 0 {
		len -= 1
	}
	var wg sync.WaitGroup
	wg.Add(len)
	users := make([]User, len)
	i, j := 0, 0
	for ; i < len; j++ {
		if ids[j] == -1 {
			continue
		}
		go func(i int, idx int64) {
			defer wg.Done()
			users[i], _ = f.GetUserByIdWithCurId(idx, userId)
		}(i, ids[i])
		i++
	}
	wg.Wait()
	// 返回关注对象列表。
	return users, nil
}

// GetFollowing 根据当前用户id来查询他的关注者列表。
func (f *FollowServiceImp) GetFollowing(userId int64) ([]User, error) {
	return getFollowing(userId)
	/*// 先查Redis，看是否有全部关注信息。
	followingIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbFollowers.SCard(middleware.Ctx, followingIdStr).Result(); 0 == cnt {
		users, _ := f.getFollowing(userId)

		go setRedisFollowing(userId, users)

		return users, nil
	}
	// Redis中有。
	UserIdStr := strconv.Itoa(int(userId))
	userIds, _ := middleware.RdbFollowing.SMembers(middleware.Ctx, UserIdStr).Result()
	len := len(userIds)
	if len > 0 {
		len -= 1
	}
	users := make([]User, len)
	wg := sync.WaitGroup{}
	wg.Add(len)
	i, j := 0, 0
	for ; i < len; j++ {
		idx, _ := strconv.Atoi(userIds[j])
		if idx == -1 {
			continue
		}
		go func(i int, idx int) {
			defer wg.Done()
			users[i], _ = f.GetUserByIdWithCurId(int64(idx), userId)
		}(i, idx)

		i++
	}
	wg.Wait()
	log.Println("从Redis中查询到所有关注者。")
	return users, nil*/
}

// 设置Redis关于所有关注的信息。
func setRedisFollowing(userId int64, users []User) {
	/*
		1-设置following_userId的所有关注id。
		2-设置following_part_id关注信息。
	*/
	// 加上-1防止脏读
	followingIdStr := strconv.Itoa(int(userId))
	redis.RdbFollowing.SAdd(redis.Ctx, followingIdStr, -1)
	// 设置过期时间
	redis.RdbFollowing.Expire(redis.Ctx, followingIdStr, config.ExpireTime)
	for i, user := range users {
		redis.RdbFollowing.SAdd(redis.Ctx, followingIdStr, user.Id)

		redis.RdbFollowingPart.SAdd(redis.Ctx, followingIdStr, user.Id)
		redis.RdbFollowingPart.SAdd(redis.Ctx, followingIdStr, -1)
		// 随机设置过期时间
		redis.RdbFollowingPart.Expire(redis.Ctx, followingIdStr, config.ExpireTime+
			time.Duration((i%10)<<8))
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

func (f *FollowServiceImp) getFollowers(userId int64) ([]User, error) {
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
	return f.getUserById(ids, userId)
}
func (f *FollowServiceImp) getUserById(ids []int64, userId int64) ([]User, error) {
	len := len(ids)
	if len > 0 {
		len -= 1
	}
	users := make([]User, len)
	var wg sync.WaitGroup
	wg.Add(len)
	i, j := 0, 0
	for ; i < len; j++ {
		// 越过-1
		if ids[j] == -1 {
			continue
		}
		//开启协程来查。
		go func(i int, idx int64) {
			defer wg.Done()
			users[i], _ = f.GetUserByIdWithCurId(idx, userId)
		}(i, ids[i])
		i++
	}
	wg.Wait()
	// 返回粉丝列表。
	return users, nil
}

// GetFollowers 根据当前用户id来查询他的粉丝列表。
func (f *FollowServiceImp) GetFollowers(userId int64) ([]User, error) {
	return getFollowers(userId)
	/*// 先查Redis，看是否有全部粉丝信息。
	followersIdStr := strconv.Itoa(int(userId))
	if cnt, _ := middleware.RdbFollowers.SCard(middleware.Ctx, followersIdStr).Result(); 0 == cnt {
		users, _ := f.getFollowers(userId)

		go setRedisFollowers(userId, users)

		return users, nil
	}
	// Redis中有。
	// 先更新有效期。
	middleware.RdbFollowers.Expire(middleware.Ctx, followersIdStr, config.ExpireTime)
	userIds, _ := middleware.RdbFollowers.SMembers(middleware.Ctx, followersIdStr).Result()
	len := len(userIds)
	if len > 0 {
		len -= 1
	}
	users := make([]User, len)
	var wg sync.WaitGroup
	wg.Add(len)
	i, j := 0, 0
	for ; i < len; j++ {
		idx, _ := strconv.Atoi(userIds[j])
		if idx == -1 {
			continue
		}
		go func(i int, idx int) {
			defer wg.Done()
			users[i], _ = f.GetUserByIdWithCurId(int64(idx), userId)
		}(i, idx)
		i++
	}
	wg.Wait()
	return users, nil*/
}

// 从数据库查所有粉丝信息。
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
		2-设置following_part_id关注信息。
	*/
	// 加上-1防止脏读。
	followersIdStr := strconv.Itoa(int(userId))
	redis.RdbFollowers.SAdd(redis.Ctx, followersIdStr, -1)
	// 设置过期时间
	redis.RdbFollowers.Expire(redis.Ctx, followersIdStr, config.ExpireTime)
	for i, user := range users {
		redis.RdbFollowers.SAdd(redis.Ctx, followersIdStr, user.Id)

		userUserIdStr := strconv.Itoa(int(user.Id))
		redis.RdbFollowingPart.SAdd(redis.Ctx, userUserIdStr, userId)
		redis.RdbFollowingPart.SAdd(redis.Ctx, userUserIdStr, -1)
		// 随机更新过期时间
		redis.RdbFollowingPart.Expire(redis.Ctx, userUserIdStr, config.ExpireTime+
			time.Duration((i%10)<<8))

		if user.IsFollow {
			redis.RdbFollowingPart.SAdd(redis.Ctx, followersIdStr, user.Id)
			redis.RdbFollowingPart.SAdd(redis.Ctx, followersIdStr, -1)
			redis.RdbFollowingPart.Expire(redis.Ctx, followersIdStr, config.ExpireTime+
				time.Duration((i%10)<<8))
		}
	}
}
