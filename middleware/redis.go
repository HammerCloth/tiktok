package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var RdbFollowers *redis.Client
var RdbFollowing *redis.Client
var RdbUser *redis.Client
var RdbFollowingPart *redis.Client

var Rdb5 *redis.Client //redis db5
var Rdb6 *redis.Client //redis db6

var RdbVCid *redis.Client  //redis db11 -- Video_id + comment_id
var RdbCInfo *redis.Client //redis db12 -- Comment_id + commentInfo

// InitRedis 初始化Redis连接。
func InitRedis() {
	RdbFollowers = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       0, // 粉丝列表信息存入 DB0.
	})
	RdbFollowing = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       1, // 关注列表信息信息存入 DB1.
	})
	RdbUser = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       2, // 关注列表和粉丝列表中的用具体信息存入 DB2.
	})
	RdbFollowingPart = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       3, // 当前用户是否关注了自己粉丝信息存入 DB1.
	})

	Rdb5 = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       5, // lls 选择将follow相关信息存入 DB5.
	})

	Rdb6 = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       6, // lls 选择将follow相关信息存入 DB6.
	})

	RdbVCid = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       11, // lsy 选择将video_id中的评论id存入 DB11.
	})

	RdbCInfo = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       12, // lsy 选择将Comment相关信息存入 DB12.
	})

}
