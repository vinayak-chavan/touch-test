package utils

import (
	"context"
	"time"

	"github.com/go-redis/redis/v8"
)

func AddToRedis(rdb *redis.Client, key string, value []byte, expiration time.Duration) error {
	redisCtx := context.Background()
	return rdb.Set(redisCtx, key, value, expiration).Err()
}

func ClearRedis(rdb *redis.Client, key string) error {
	redisCtx := context.Background()
	return rdb.Del(redisCtx, key).Err()
}
