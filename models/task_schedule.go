package models

import (
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type TaskSchedule struct {
	BaseModel
	ScheduleType utils.CronJobScheduleType
	DirectoryID  uuid.UUID
	Directory    Directory `gorm:"foreignKey:DirectoryID"`
}

func DeleteTaskSchedule(ctx *gin.Context, db *gorm.DB, id uuid.UUID) error {
	var taskSchedule *TaskSchedule
	result := db.WithContext(ctx).Model(&TaskSchedule{}).Where("id = ?", id).Delete(&taskSchedule)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}

func CreateTaskSchedule(ctx *gin.Context, db *gorm.DB, directory Directory, scheduleType utils.CronJobScheduleType) (*TaskSchedule, error) {
	taskSchedule := &TaskSchedule{
		DirectoryID:  directory.ID,
		ScheduleType: scheduleType,
	}

	result := db.Debug().WithContext(ctx).Model(&TaskSchedule{}).Create(&taskSchedule)
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

func GetTaskSchedules(ctx *gin.Context, db *gorm.DB) ([]TaskSchedule, error) {
	var tasks []TaskSchedule

	result := db.WithContext(ctx).Model(&TaskSchedule{}).Find(&tasks)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return tasks, nil
}

func GetTaskSchedulesByScheduleType(ctx *gin.Context, db *gorm.DB, scheduleType utils.CronJobScheduleType) ([]TaskSchedule, error) {
	var tasks []TaskSchedule

	result := db.WithContext(ctx).Model(&TaskSchedule{}).Where("schedule_type = ?", scheduleType).Find(&tasks)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return tasks, nil
}

func GetTaskSchedule(ctx *gin.Context, db *gorm.DB, id uuid.UUID) (*TaskSchedule, error) {
	var task *TaskSchedule

	result := db.WithContext(ctx).Model(&TaskSchedule{}).Where("id = ?", id).Find(&task)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return task, nil
}
