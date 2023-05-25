package clients

import (
	redis "github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var Clients *GlobalClients

type GlobalClients struct {
	DB    *gorm.DB
	Redis *redis.Client
	Cron  *Cron
}

func NewGlobalClients(
	postgres *gorm.DB,
	redis *redis.Client,
	cron *Cron,
) *GlobalClients {
	return &GlobalClients{
		DB:    postgres,
		Redis: redis,
		Cron:  cron,
	}
}
