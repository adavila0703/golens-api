package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IgnoredDirectories struct {
	BaseModel
	DirectoryName string
}

func (i *IgnoredDirectories) Test() {

}

func CreateIgnoredDirectory(ctx *gin.Context, db *gorm.DB, directoryName string) error {
	ignoredDirectories := &IgnoredDirectories{
		DirectoryName: directoryName,
	}

	results := db.WithContext(ctx).Create(&ignoredDirectories)

	if err := results.Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetIgnoredDirectories(ctx *gin.Context, db *gorm.DB) []IgnoredDirectories {
	var ignoredDirectories []IgnoredDirectories
	results := db.WithContext(ctx).Model(IgnoredDirectories{}).Find(&ignoredDirectories)

	if results.RowsAffected == 0 {
		return nil
	}

	return ignoredDirectories
}

func DeleteIgnoredDirectory(ctx *gin.Context, db *gorm.DB, id uuid.UUID) error {
	var ignoredDirectory *IgnoredDirectories
	result := db.WithContext(ctx).Model(&IgnoredDirectories{}).Where("id = ?", id).Delete(&ignoredDirectory)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}
