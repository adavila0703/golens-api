package clients

import (
	"golens-api/coverage"

	"gorm.io/gorm"
)

var Clients *GlobalClients

type GlobalClients struct {
	DB   *gorm.DB
	Cron ICron
	Cov  coverage.ICoverage
}

func NewGlobalClients(
	postgres *gorm.DB,
	cron ICron,
	cov coverage.ICoverage,
) *GlobalClients {
	return &GlobalClients{
		DB:   postgres,
		Cron: cron,
		Cov:  cov,
	}
}
