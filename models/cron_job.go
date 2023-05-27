package models

import (
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type CronJob struct {
	BaseModel
	Schedule     string
	ScheduleType utils.CronJobScheduleType
	EntryID      cron.EntryID
}

func UpdateCronJob(ctx *gin.Context, db *gorm.DB, cronJob *CronJob) error {
	result := db.WithContext(ctx).Updates(&cronJob)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}

func CreateCronJob(ctx *gin.Context, db *gorm.DB, scheduleType utils.CronJobScheduleType, entryID cron.EntryID) (*CronJob, error) {
	schedule := utils.GetCronSchedule(scheduleType)
	cronJob := &CronJob{
		Schedule:     schedule,
		ScheduleType: scheduleType,
		EntryID:      entryID,
	}

	result := db.WithContext(ctx).Model(&CronJob{}).Where(&cronJob).FirstOrCreate(&cronJob)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return cronJob, nil
}

func CronJobExists(ctx *gin.Context, db *gorm.DB, schedule utils.CronJobScheduleType) bool {
	var cronJob *CronJob
	result := db.WithContext(ctx).Model(&CronJob{}).Where("schedule_type = ?", schedule).Find(&cronJob)

	return result.RowsAffected > 0
}

// there should never be more than one of these cron jobs
func GetCronJob(ctx *gin.Context, db *gorm.DB, scheduleType utils.CronJobScheduleType) (*CronJob, error) {
	var cronJob *CronJob

	result := db.WithContext(ctx).Model(&CronJob{}).Where("schedule_type = ?", scheduleType).Find(&cronJob)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return cronJob, nil
}

func DeleteCronJob(ctx *gin.Context, db *gorm.DB, id uuid.UUID) error {
	var cronJob *CronJob
	result := db.WithContext(ctx).Model(&CronJob{}).Where("id = ?", id).Delete(&cronJob)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}
