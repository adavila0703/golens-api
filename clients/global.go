package clients

import (
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var Clients *GlobalClients

type GlobalClients struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewGlobalClients(
	postgres *gorm.DB,
	redis *redis.Client,
) *GlobalClients {
	return &GlobalClients{
		DB:    postgres,
		Redis: redis,
	}
}
