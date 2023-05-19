package clients

import (
	redis "github.com/redis/go-redis/v9"
)

func NewRedisClient(options *redis.Options) *redis.Client {
	return redis.NewClient(options)
}
