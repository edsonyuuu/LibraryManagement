package core

import (
	"LibraryManagementV1/LM_V4/global"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
)

func InitRedis() {
	// 创建 Redis 客户端连接
	global.RedisConn = redis.NewClient(&redis.Options{
		Addr:     global.Config.Redis.Addr(),
		Password: global.Config.Redis.Password, // Redis 未设置密码时为空
		DB:       1,
	})
	// 测试连接是否成功
	_, err := global.RedisConn.Ping(context.Background()).Result()
	if err != nil {
		global.Log.Error("Failed to connect to Redis for InfoCache: %v", err)
		return
	}
	fmt.Println("Connected to Redis")
}
