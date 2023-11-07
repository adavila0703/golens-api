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
	Name string
	Type IgnoreType
}

func CreateIgnored(ctx *gin.Context, db *gorm.DB, directoryName string, ignoredType IgnoreType) error {
	ignoredDirectories := &Ignored{
		Name: directoryName,
		Type: ignoredType,
	}

	results := db.WithContext(ctx).Create(&ignoredDirectories)

	if err := results.Error; err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func GetIgnored(ctx *gin.Context, db *gorm.DB, ignoredType IgnoreType) []Ignored {
	var ignoredDirectories []Ignored
	results := db.WithContext(ctx).Model(Ignored{Type: ignoredType}).Find(&ignoredDirectories)

	if results.RowsAffected == 0 {
		return nil
	}

	return ignoredDirectories
}

func DeleteIgnored(ctx *gin.Context, db *gorm.DB, id uuid.UUID) error {
	var ignoredDirectory *Ignored
	result := db.WithContext(ctx).Model(&Ignored{}).Where("id = ?", id).Delete(&ignoredDirectory)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}
