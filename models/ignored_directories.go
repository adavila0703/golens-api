package models

import (
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IgnoredDirectories struct {
	BaseModel
	DirectoryName string
}

func AddIgnoredDirectory(ctx *gin.Context, db *gorm.DB, directoryName string) error {
	ignoredDirectories := &IgnoredDirectories{
		DirectoryName: directoryName,
	}

	results := db.WithContext(ctx).Create(&ignoredDirectories)

	if err := results.Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetIgnoredDirectories(ctx *gin.Context, db *gorm.DB) []string {
	var ignoredDirectories []IgnoredDirectories
	results := db.WithContext(ctx).Model(IgnoredDirectories{}).First(&ignoredDirectories)

	if results.RowsAffected == 0 {
		return nil
	}

	var directoryNames []string

	for _, directory := range ignoredDirectories {
		directoryNames = append(directoryNames, directory.DirectoryName)
	}

	return directoryNames
}
