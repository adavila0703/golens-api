package models

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TaskSchedule struct {
	BaseModel
	TaskID      int
	Schedule    string
	DirectoryID uuid.UUID
	Directory   Directory `gorm:"foreignKey:DirectoryID"`
}

func CreateTaskSchedule(ctx *gin.Context, db *gorm.DB, directory *Directory) (*TaskSchedule, error) {
	taskSchedule := &TaskSchedule{
		DirectoryID: directory.ID,
	}

	result := db.WithContext(ctx).Model(&TaskSchedule{}).Where(&taskSchedule).FirstOrCreate(&taskSchedule)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return taskSchedule, nil
}

func GetTaskScheduleByDirectoryID(ctx *gin.Context, db *gorm.DB, directoryID uuid.UUID) (*TaskSchedule, error) {
	taskSchedule := &TaskSchedule{
		DirectoryID: directoryID,
	}

	result := db.WithContext(ctx).Model(&TaskSchedule{}).Where(&taskSchedule).Find(&taskSchedule)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return taskSchedule, nil
}

func GetTaskSchedules(ctx *gin.Context, db *gorm.DB) ([]*TaskSchedule, error) {
	var tasks []*TaskSchedule

	result := db.WithContext(ctx).Model(&TaskSchedule{}).Find(&tasks)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return tasks, nil
}
