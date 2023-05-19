package models

import (
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Directory struct {
	BaseModel
	Path         string
	CoverageName string
}

func CreateDirectory(ctx *gin.Context, db *gorm.DB, path string) error {
	coverageName := utils.GetCoverageNameFromPath(path)
	directory := &Directory{
		Path:         path,
		CoverageName: coverageName,
	}

	result := db.WithContext(ctx).Model(&Directory{}).Where(directory).FirstOrCreate(&directory)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}

func GetDirectories(ctx *gin.Context, db *gorm.DB) ([]Directory, error) {
	var directories []Directory

	result := db.WithContext(ctx).Model(&Directory{}).Find(&directories)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return directories, nil
}
