package models

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Setting struct {
	BaseModel
	IgnoredDirectories []string
}

func AddIgnoredDirectory(ctx *gin.Context, db *gorm.DB, directoryName string) error {
	var settings Setting
	results := db.WithContext(ctx).Model(Setting{}).First(&settings)

	if err := results.Error; err != nil {
		return errors.WithStack(err)
	}

	settings.IgnoredDirectories = append(settings.IgnoredDirectories, directoryName)

	results = db.WithContext(ctx).Model(Setting{}).Updates(&settings)

	if err := results.Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetIgnoredDirectories(ctx *gin.Context, db *gorm.DB) ([]string, error) {
	var settings Setting
	results := db.WithContext(ctx).Model(Setting{}).First(&settings)

	if err := results.Error; err != nil {
		return nil, errors.WithStack(err)
	}

	return settings.IgnoredDirectories, nil
}
