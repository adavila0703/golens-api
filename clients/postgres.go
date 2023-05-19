package clients

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func NewPostgresClient(dns string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{
		// shutting off logging since its too noisy
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return db, nil
}

func NewPostgresClientMock() (*gorm.DB, sqlmock.Sqlmock, func() error, error) {
	db, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	dialector := postgres.New(postgres.Config{
		DSN:                  "sqlmock_db_0",
		DriverName:           "postgres",
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	mockDB, err := gorm.Open(dialector, &gorm.Config{
		Logger: logger.Default.LogMode(logger.Error),
	})
	if err != nil {
		return nil, nil, nil, errors.WithStack(err)
	}

	return mockDB, mock, db.Close, nil
}
