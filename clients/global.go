package clients

import (
	"gorm.io/gorm"
)

var Clients *GlobalClients

type GlobalClients struct {
	DB   *gorm.DB
	Cron *Cron
}

func NewGlobalClients(
	postgres *gorm.DB,
	cron *Cron,
) *GlobalClients {
	return &GlobalClients{
		DB:   postgres,
		Cron: cron,
	}
}
