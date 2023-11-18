package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type IgnoreType int

const (
	None IgnoreType = iota
	DirectoryType
	PathType
	FileType
	PackageType
)

type Ignored struct {
	BaseModel
	DirectoryName string
	Name          string
	Type          IgnoreType
}

func CreateIgnored(
	ctx *gin.Context,
	db *gorm.DB,
	directoryName string,
	name string,
	ignoredType IgnoreType,
) error {
	ignored := &Ignored{
		DirectoryName: directoryName,
		Name:          name,
		Type:          ignoredType,
	}

	results := db.WithContext(ctx).Create(&ignored)

	if err := results.Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetIgnored(ctx *gin.Context, db *gorm.DB, ignoredType IgnoreType) []Ignored {
	var ignored []Ignored
	results := db.
		WithContext(ctx).
		Model(&Ignored{}).
		Where("type = ?", ignoredType).
		Find(&ignored)

	if results.RowsAffected == 0 {
		return nil
	}

	return ignored
}

func DeleteIgnored(ctx *gin.Context, db *gorm.DB, id uuid.UUID) error {
	var ignored *Ignored
	result := db.WithContext(ctx).Model(&Ignored{}).Where("id = ?", id).Delete(&ignored)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}
