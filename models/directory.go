package models

import (
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Directory struct {
	BaseModel
	Path         string
	CoverageName string
}

func DirectoryExists(ctx *gin.Context, db *gorm.DB, path string) (bool, error) {
	var directory *Directory

	result := db.WithContext(ctx).Model(&Directory{}).Where("path = ?", path).Find(&directory)

	if result.Error != nil {
		return false, errors.WithStack(result.Error)
	}

	found := result.RowsAffected > 0

	return found, nil
}

func CreateDirectory(ctx *gin.Context, db *gorm.DB, path string) (*Directory, error) {
	coverageName := utils.GetCoverageNameFromPath(path)
	directory := &Directory{
		Path:         path,
		CoverageName: coverageName,
	}

	result := db.WithContext(ctx).Model(&Directory{}).Where(directory).FirstOrCreate(&directory)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return directory, nil
}

func GetDirectories(ctx *gin.Context, db *gorm.DB) ([]Directory, error) {
	var directories []Directory

	result := db.WithContext(ctx).Model(&Directory{}).Find(&directories)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return directories, nil
}

func GetDirectory(ctx *gin.Context, db *gorm.DB, id uuid.UUID) (*Directory, bool, error) {
	var directory *Directory

	result := db.WithContext(ctx).Model(&Directory{}).Where("id = ?", id).Find(&directory)

	if result.Error != nil {
		return nil, false, errors.WithStack(result.Error)
	}

	found := result.RowsAffected > 0

	return directory, found, nil
}

func DeleteDirectory(ctx *gin.Context, db *gorm.DB, id uuid.UUID) error {
	var directory *Directory
	result := db.WithContext(ctx).Model(&Directory{}).Where("id = ?", id).Delete(&directory)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}
