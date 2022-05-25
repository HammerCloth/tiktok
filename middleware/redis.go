package middleware

import (
	"context"
	"github.com/go-redis/redis/v8"
)

var Ctx = context.Background()
var Rdb *redis.Client

// InitRedis 初始化Redis连接。
func InitRedis() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     "106.14.75.229:6379",
		Password: "tiktok",
		DB:       0, // lls 选择将follow相关信息存入 DB0.
	})
}
