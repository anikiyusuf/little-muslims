package cache

import "github.com/redis/go-redis/v9"


type RedisCache struct {
	*redis.Client
}

func NewRedisCache(address string, pass string) *RedisCache {
	return &RedisCache{redis.NewClient(&redis.Options{
		Addr: address,
		Password: pass,
	})}
}
