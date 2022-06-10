package service

import (
	"TikTok/config"
	"TikTok/dao"
	"TikTok/middleware/rabbitmq"
	"TikTok/middleware/redis"
	"errors"
	"log"
	"strconv"
	"strings"
	"sync"
	"time"
)

type LikeServiceImpl struct {
	VideoService
	UserService
}

//IsFavourite 根据userId,videoId查询点赞状态 这边可以快一点,通过查询两个Redis DB;
//step1：查询Redis LikeUserId(key:strUserId)是否已经加载过此信息，通过是否存在value:videoId 判断点赞状态;
//step2:如LikeUserId没有对应信息，查询LikeVideoId(key：strVideoId)是否已经加载过此信息，通过是否存在value:userId 判断点赞状态;
//step3:LikeUserId LikeVideoId中都没有对应key,维护LikeUserId对应key，并通过查询key：strUserId中是否存在value:videoId 判断点赞状态;
func (like *LikeServiceImpl) IsFavourite(videoId int64, userId int64) (bool, error) {
	//将int64 userId转换为 string strUserId
	strUserId := strconv.FormatInt(userId, 10)
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	//step1:查询Redis LikeUserId,key：strUserId中是否存在value:videoId,key中存在value 返回true，不存在返回false
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("方法:IsFavourite RedisLikeUserId query key失败：%v", err)
			return false, err
		}
		exist, err1 := redis.RdbLikeUserId.SIsMember(redis.Ctx, strUserId, videoId).Result()
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err1 != nil {
			log.Printf("方法:IsFavourite RedisLikeUserId query value失败：%v", err1)
			return false, err1
		}
		log.Printf("方法:IsFavourite RedisLikeUserId query value成功")
		return exist, nil
	} else { //step2:LikeUserId不存在key,查询Redis LikeVideoId,key中存在value 返回true，不存在返回false
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回默认false,返回错误信息
			if err != nil {
				log.Printf("方法:IsFavourite RedisLikeVideoId query key失败：%v", err)
				return false, err
			}
			exist, err1 := redis.RdbLikeVideoId.SIsMember(redis.Ctx, strVideoId, userId).Result()
			//如果有问题，说明查询redis失败,返回默认false,返回错误信息
			if err1 != nil {
				log.Printf("方法:IsFavourite RedisLikeVideoId query value失败：%v", err1)
				return false, err1
			}
			log.Printf("方法:IsFavourite RedisLikeVideoId query value成功")
			return exist, nil
		} else {
			//key:strUserId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("方法:IsFavourite RedisLikeUserId add value失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return false, err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("方法:IsFavourite RedisLikeUserId 设置有效期失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return false, err
			}
			//step3:LikeUserId LikeVideoId中都没有对应key，通过userId查询likes表,返回所有点赞videoId，并维护到Redis LikeUserId(key:strUserId)
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			//如果有问题，说明查询数据库失败，返回默认false,返回错误信息："get likeVideoIdList failed"
			if err1 != nil {
				log.Printf(err1.Error())
				return false, err1
			}
			//维护Redis LikeUserId(key:strUserId)，遍历videoIdList加入
			for _, likeVideoId := range videoIdList {
				redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId)
			}
			//查询Redis LikeUserId,key：strUserId中是否存在value:videoId,存在返回true,不存在返回false
			exist, err2 := redis.RdbLikeUserId.SIsMember(redis.Ctx, strUserId, videoId).Result()
			//如果有问题，说明操作redis失败,返回默认false,返回错误信息
			if err2 != nil {
				log.Printf("方法:IsFavourite RedisLikeUserId query value失败：%v", err2)
				return false, err2
			}
			log.Printf("方法:IsFavourite RedisLikeUserId query value成功")
			return exist, nil
		}
	}
}

//FavouriteCount 根据videoId获取对应点赞数量;
//step1：查询Redis LikeVideoId(key:strVideoId)是否已经加载过此信息，通过set集合中userId个数，获取点赞数量;
//step2：LikeVideoId中都没有对应key，维护LikeVideoId对应key，再通过set集合中userId个数，获取点赞数量;
func (like *LikeServiceImpl) FavouriteCount(videoId int64) (int64, error) {
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	//step1 如果key:strVideoId存在 则计算集合中userId个数
	if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认false,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId query key失败：%v", err)
			return 0, err
		}
		//获取集合中userId个数
		count, err1 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err1 != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId query count 失败：%v", err1)
			return 0, err1
		}
		log.Printf("方法:FavouriteCount RedisLikeVideoId query count 成功")
		return count - 1, nil //去掉DefaultRedisValue
	} else {
		//key:strVideoId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
		if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId add value失败")
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		//给键值设置有效期，类似于gc机制
		_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId,
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId 设置有效期失败")
			redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
			return 0, err
		}
		//如果Redis LikeVideoId不存在此key,通过videoId查询likes表,返回所有点赞userId，并维护到Redis LikeVideoId(key:strVideoId)
		//再通过set集合中userId个数,获取点赞数量
		userIdList, err1 := dao.GetLikeUserIdList(videoId)
		//如果有问题，说明查询数据库失败，返回默认0,返回错误信息："get likeUserIdList failed"
		if err1 != nil {
			log.Printf(err1.Error())
			return 0, err1
		}
		//维护Redis LikeVideoId(key:strVideoId)，遍历userIdList加入
		for _, likeUserId := range userIdList {
			redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId)
		}
		//再通过set集合中userId个数,获取点赞数量
		count, err2 := redis.RdbLikeVideoId.SCard(redis.Ctx, strVideoId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err2 != nil {
			log.Printf("方法:FavouriteCount RedisLikeVideoId query count 失败：%v", err2)
			return 0, err2
		}
		log.Printf("方法:FavouriteCount RedisLikeVideoId query count 成功")
		return count - 1, nil //去掉DefaultRedisValue
	}
}

// FavouriteAction 根据userId，videoId,actionType对视频进行点赞或者取消赞操作;
//step1: 维护Redis LikeUserId(key:strUserId),添加或者删除value:videoId,LikeVideoId(key:strVideoId),添加或者删除value:userId;
//step2：更新数据库likes表;
func (like *LikeServiceImpl) FavouriteAction(userId int64, videoId int64, actionType int32) error {
	//将int64 videoId转换为 string strVideoId
	strUserId := strconv.FormatInt(userId, 10)
	//将int64 videoId转换为 string strVideoId
	strVideoId := strconv.FormatInt(videoId, 10)
	//将要操作数据库likes表的信息打入消息队列RmqLikeAdd或者RmqLikeDel
	//拼接打入信息
	sb := strings.Builder{}
	sb.WriteString(strUserId)
	sb.WriteString(" ")
	sb.WriteString(strVideoId)

	//step1:维护Redis LikeUserId、LikeVideoId;
	//执行点赞操作维护
	if actionType == config.LikeAction {
		//查询Redis LikeUserId(key:strUserId)是否已经加载过此信息
		if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回错误信息
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId query key失败：%v", err)
				return err
			} //如果加载过此信息key:strUserId，则加入value:videoId
			//如果redis LikeUserId 添加失败，数据库操作成功，会有脏数据，所以只有redis操作成功才执行数据库likes表操作
			if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId add value失败：%v", err1)
				return err1
			} else {
				//如果数据库操作失败了，redis是正确数据，客户端显示的是点赞成功，不会影响后续结果
				//只有当该用户取消所有点赞视频的时候redis才会重新加载数据库信息，这时候因为取消赞了必然和数据库信息一致
				//同样这条信息消费成功与否也不重要，因为redis是正确信息,理由如上
				rabbitmq.RmqLikeAdd.Publish(sb.String())
			}
		} else { //如果不存在，则维护Redis LikeUserId 新建key:strUserId,设置过期时间，加入DefaultRedisValue，
			//key:strUserId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			//通过userId查询likes表,返回所有点赞videoId，加入key:strUserId集合中,
			//再加入当前videoId,再更新likes表此条数据
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId add value失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId 设置有效期失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			//如果有问题，说明查询失败，返回错误信息："get likeVideoIdList failed"
			if err1 != nil {
				return err1
			}
			//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql数据一致性
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					log.Printf("方法:FavouriteAction RedisLikeUserId add value失败")
					redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId add value失败：%v", err2)
				return err2
			} else {
				rabbitmq.RmqLikeAdd.Publish(sb.String())
			}
		}
		//查询Redis LikeVideoId(key:strVideoId)是否已经加载过此信息
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回错误信息
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId query key失败：%v", err)
				return err
			} //如果加载过此信息key:strVideoId，则加入value:userId
			//如果redis LikeVideoId 添加失败，返回错误信息
			if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败：%v", err1)
				return err1
			}
		} else { //如果不存在，则维护Redis LikeVideoId 新建key:strVideoId，设置有效期，加入DefaultRedisValue
			//通过videoId查询likes表,返回所有点赞userId，加入key:strVideoId集合中,
			//再加入当前userId,再更新likes表此条数据
			//key:strVideoId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败")
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId 设置有效期失败")
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			userIdList, err1 := dao.GetLikeUserIdList(videoId)
			//如果有问题，说明查询失败，返回错误信息："get likeUserIdList failed"
			if err1 != nil {
				return err1
			}
			//遍历userIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql数据一致性
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败")
					redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败：%v", err2)
				return err2
			}
		}
	} else { //执行取消赞操作维护
		//查询Redis LikeUserId(key:strUserId)是否已经加载过此信息
		if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回错误信息
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId query key失败：%v", err)
				return err
			} //防止出现redis数据不一致情况，当redis删除操作成功，才执行数据库更新操作
			if _, err1 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err1 != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId del value失败：%v", err1)
				return err1
			} else {
				//后续数据库的操作，可以在mq里设置若执行数据库更新操作失败，重新消费该信息
				rabbitmq.RmqLikeDel.Publish(sb.String())
			}
		} else { //如果不存在，则维护Redis LikeUserId 新建key:strUserId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库
			// 还没更新完出现脏读，或者数据库操作失败造成的脏读
			//通过userId查询likes表,返回所有点赞videoId，加入key:strUserId集合中,
			//再删除当前videoId,再更新likes表此条数据
			//key:strUserId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId add value失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId 设置有效期失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return err
			}
			videoIdList, err1 := dao.GetLikeVideoIdList(userId)
			//如果有问题，说明查询失败，返回错误信息："get likeVideoIdList failed"
			if err1 != nil {
				return err1
			}
			//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql 数据原子性
			for _, likeVideoId := range videoIdList {
				if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
					log.Printf("方法:FavouriteAction RedisLikeUserId add value失败")
					redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbLikeUserId.SRem(redis.Ctx, strUserId, videoId).Result(); err2 != nil {
				log.Printf("方法:FavouriteAction RedisLikeUserId del value失败：%v", err2)
				return err2
			} else {
				rabbitmq.RmqLikeDel.Publish(sb.String())
			}
		}

		//查询Redis LikeVideoId(key:strVideoId)是否已经加载过此信息
		if n, err := redis.RdbLikeVideoId.Exists(redis.Ctx, strVideoId).Result(); n > 0 {
			//如果有问题，说明查询redis失败,返回错误信息
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId query key失败：%v", err)
				return err
			} //如果加载过此信息key:strVideoId，则删除value:userId
			//如果redis LikeVideoId 删除失败，返回错误信息
			if _, err1 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err1 != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId del value失败：%v", err1)
				return err1
			}
		} else { //如果不存在，则维护Redis LikeVideoId 新建key:strVideoId,加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库
			// 还没更新完出现脏读，或者数据库操作失败造成的脏读
			//通过videoId查询likes表,返回所有点赞userId，加入key:strVideoId集合中,
			//再删除当前userId,再更新likes表此条数据
			//key:strVideoId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
			if _, err := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, config.DefaultRedisValue).Result(); err != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败")
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}
			//给键值设置有效期，类似于gc机制
			_, err := redis.RdbLikeVideoId.Expire(redis.Ctx, strVideoId,
				time.Duration(config.OneMonth)*time.Second).Result()
			if err != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId 设置有效期失败")
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err
			}

			userIdList, err1 := dao.GetLikeUserIdList(videoId)
			//如果有问题，说明查询失败，返回错误信息："get likeUserIdList failed"
			if err1 != nil {
				redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
				return err1
			}
			//遍历userIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
			//保证redis与mysql数据一致性
			for _, likeUserId := range userIdList {
				if _, err1 := redis.RdbLikeVideoId.SAdd(redis.Ctx, strVideoId, likeUserId).Result(); err1 != nil {
					log.Printf("方法:FavouriteAction RedisLikeVideoId add value失败")
					redis.RdbLikeVideoId.Del(redis.Ctx, strVideoId)
					return err1
				}
			}
			//这样操作理由同上
			if _, err2 := redis.RdbLikeVideoId.SRem(redis.Ctx, strVideoId, userId).Result(); err2 != nil {
				log.Printf("方法:FavouriteAction RedisLikeVideoId del value失败：%v", err2)
				return err2
			}
		}
	}
	return nil
}

//GetFavouriteList 根据userId，curId(当前用户Id),返回userId的点赞列表;
//step1：查询Redis LikeUserId(key:strUserId)是否已经加载过此信息，获取集合中全部videoId，并添加到点赞列表集合中;
//step2：LikeUserId中都没有对应key，维护LikeUserId对应key，同时添加到点赞列表集合中;
func (like *LikeServiceImpl) GetFavouriteList(userId int64, curId int64) ([]Video, error) {
	//将int64 userId转换为 string strUserId
	strUserId := strconv.FormatInt(userId, 10)
	//step1:查询Redis LikeUserId,如果key：strUserId存在,则获取集合中全部videoId
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err != nil {
			log.Printf("方法:GetFavouriteList RedisLikeVideoId query key失败：%v", err)
			return nil, err
		}
		//获取集合中全部videoId
		videoIdList, err1 := redis.RdbLikeUserId.SMembers(redis.Ctx, strUserId).Result()
		//如果有问题，说明查询redis失败,返回默认nil,返回错误信息
		if err1 != nil {
			log.Printf("方法:GetFavouriteList RedisLikeVideoId get values失败：%v", err1)
			return nil, err1
		}
		//提前开辟点赞列表空间
		favoriteVideoList := new([]Video)
		//采用协程并发将Video类型对象添加到集合中去
		i := len(videoIdList) - 1 //去掉DefaultRedisValue
		if i == 0 {
			return *favoriteVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			//将string videoId转换为 int64 VideoId
			videoId, _ := strconv.ParseInt(videoIdList[j], 10, 64)
			if videoId == config.DefaultRedisValue {
				continue
			}
			go like.addFavouriteVideoList(videoId, curId, favoriteVideoList, &wg)
		}
		wg.Wait()
		return *favoriteVideoList, nil
	} else { //如果Redis LikeUserId不存在此key,通过userId查询likes表,返回所有点赞videoId，并维护到Redis LikeUserId(key:strUserId)
		//key:strUserId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
		if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
			log.Printf("方法:GetFavouriteList RedisLikeUserId add value失败")
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err
		}
		//给键值设置有效期，类似于gc机制
		_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			log.Printf("方法:GetFavouriteList RedisLikeUserId 设置有效期失败")
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err
		}
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		//如果有问题，说明查询数据库失败，返回nil和错误信息:"get likeVideoIdList failed"
		if err1 != nil {
			log.Println(err1.Error())
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return nil, err1
		}
		//遍历videoIdList,添加进key的集合中，若失败，删除key，并返回错误信息，这么做的原因是防止脏读，
		//保证redis与mysql数据一致性
		for _, likeVideoId := range videoIdList {
			if _, err2 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err2 != nil {
				log.Printf("方法:GetFavouriteList RedisLikeUserId add value失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return nil, err2
			}
		}
		//提前开辟点赞列表空间
		favoriteVideoList := new([]Video)
		//采用协程并发将Video类型对象添加到集合中去
		i := len(videoIdList) - 1 //去掉DefaultRedisValue
		if i == 0 {
			return *favoriteVideoList, nil
		}
		var wg sync.WaitGroup
		wg.Add(i)
		for j := 0; j <= i; j++ {
			if videoIdList[j] == config.DefaultRedisValue {
				continue
			}
			go like.addFavouriteVideoList(videoIdList[j], curId, favoriteVideoList, &wg)
		}
		wg.Wait()
		return *favoriteVideoList, nil
	}
}

//addFavouriteVideoList 根据videoId,登录用户curId，添加视频对象到点赞列表空间
func (like *LikeServiceImpl) addFavouriteVideoList(videoId int64, curId int64, favoriteVideoList *[]Video, wg *sync.WaitGroup) {
	defer wg.Done()
	//调用videoService接口，GetVideo：根据videoId，当前用户id:curId，返回Video类型对象
	video, err := like.GetVideo(videoId, curId)
	//如果没有获取这个video_id的视频，视频可能被删除了,打印异常,并且不加入此视频
	if err != nil {
		log.Println(errors.New("this favourite video is miss"))
		return
	}
	//将Video类型对象添加到集合中去
	*favoriteVideoList = append(*favoriteVideoList, video)
}

//TotalFavourite 根据userId获取这个用户总共被点赞数量
func (like *LikeServiceImpl) TotalFavourite(userId int64) (int64, error) {
	//根据userId获取这个用户的发布视频列表信息
	videoIdList, err := like.GetVideoIdList(userId)
	if err != nil {
		log.Printf(err.Error())
		return 0, err
	}
	var sum int64 //该用户的总被点赞数
	//提前开辟空间,存取每个视频的点赞数
	videoLikeCountList := new([]int64)
	//采用协程并发将对应videoId的点赞数添加到集合中去
	i := len(videoIdList)
	var wg sync.WaitGroup
	wg.Add(i)
	for j := 0; j < i; j++ {
		go like.addVideoLikeCount(videoIdList[j], videoLikeCountList, &wg)
	}
	wg.Wait()
	//遍历累加，求总被点赞数
	for _, count := range *videoLikeCountList {
		sum += count
	}
	return sum, nil
}

//FavouriteVideoCount 根据userId获取这个用户点赞视频数量
func (like *LikeServiceImpl) FavouriteVideoCount(userId int64) (int64, error) {
	//将int64 userId转换为 string strUserId
	strUserId := strconv.FormatInt(userId, 10)
	//step1:查询Redis LikeUserId,如果key：strUserId存在,则获取集合中元素个数
	if n, err := redis.RdbLikeUserId.Exists(redis.Ctx, strUserId).Result(); n > 0 {
		//如果有问题，说明查询redis失败,返回默认0,返回错误信息
		if err != nil {
			log.Printf("方法:FavouriteVideoCount RdbLikeUserId query key失败：%v", err)
			return 0, err
		} else {
			count, err1 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
			//如果有问题，说明操作redis失败,返回默认0,返回错误信息
			if err1 != nil {
				log.Printf("方法:FavouriteVideoCount RdbLikeUserId query count 失败：%v", err1)
				return 0, err1
			}
			log.Printf("方法:FavouriteVideoCount RdbLikeUserId query count 成功")
			return count - 1, nil //去掉DefaultRedisValue

		}
	} else { //如果Redis LikeUserId不存在此key,通过userId查询likes表,返回所有点赞videoId，并维护到Redis LikeUserId(key:strUserId)
		//再通过set集合中userId个数,获取点赞数量
		//key:strUserId，加入value:DefaultRedisValue,过期才会删，防止删最后一个数据的时候数据库还没更新完出现脏读，或者数据库操作失败造成的脏读
		if _, err := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, config.DefaultRedisValue).Result(); err != nil {
			log.Printf("方法:FavouriteVideoCount RedisLikeUserId add value失败")
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return 0, err
		}
		//给键值设置有效期，类似于gc机制
		_, err := redis.RdbLikeUserId.Expire(redis.Ctx, strUserId,
			time.Duration(config.OneMonth)*time.Second).Result()
		if err != nil {
			log.Printf("方法:FavouriteVideoCount RedisLikeUserId 设置有效期失败")
			redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
			return 0, err
		}
		videoIdList, err1 := dao.GetLikeVideoIdList(userId)
		//如果有问题，说明查询数据库失败，返回默认0,返回错误信息："get likeVideoIdList failed"
		if err1 != nil {
			log.Printf(err1.Error())
			return 0, err1
		}
		//维护Redis LikeUserId(key:strUserId)，遍历videoIdList加入
		for _, likeVideoId := range videoIdList {
			if _, err1 := redis.RdbLikeUserId.SAdd(redis.Ctx, strUserId, likeVideoId).Result(); err1 != nil {
				log.Printf("方法:FavouriteVideoCount RedisLikeUserId add value失败")
				redis.RdbLikeUserId.Del(redis.Ctx, strUserId)
				return 0, err1
			}
		}
		//再通过set集合中videoId个数,获取点赞数量
		count, err2 := redis.RdbLikeUserId.SCard(redis.Ctx, strUserId).Result()
		//如果有问题，说明操作redis失败,返回默认0,返回错误信息
		if err2 != nil {
			log.Printf("方法:FavouriteVideoCount RdbLikeUserId query count 失败：%v", err2)
			return 0, err2
		}
		log.Printf("方法:FavouriteVideoCount RdbLikeUserId query count 成功")
		return count - 1, nil //去掉DefaultRedisValue
	}
}

//addVideoLikeCount 根据videoId，将该视频点赞数加入对应提前开辟好的空间内
func (like *LikeServiceImpl) addVideoLikeCount(videoId int64, videoLikeCountList *[]int64, wg *sync.WaitGroup) {
	defer wg.Done()
	//调用FavouriteCount：根据videoId,获取点赞数
	count, err := like.FavouriteCount(videoId)
	if err != nil {
		//如果有错误，输出错误信息，并不加入该视频点赞数
		log.Printf(err.Error())
		return
	}
	*videoLikeCountList = append(*videoLikeCountList, count)
}

//GetLikeService 解决likeService调videoService,videoService调userService,useService调likeService循环依赖的问题
func GetLikeService() LikeServiceImpl {
	var userService UserServiceImpl
	var videoService VideoServiceImpl
	var likeService LikeServiceImpl
	userService.LikeService = &likeService
	likeService.VideoService = &videoService
	videoService.UserService = &userService
	return likeService
}
