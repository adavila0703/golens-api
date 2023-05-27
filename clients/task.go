package clients

import (
	"fmt"
	"golens-api/models"
	"golens-api/utils"
	"log"

	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

func GetUpdateTaskFunc(cronSchedule utils.CronJobScheduleType) func() {
	switch cronSchedule {
	case utils.EveryMinute:
		return UpdateCoverageTask_EveryMinute
	case utils.EveryDayAt12AM:
		return UpdateCoverageTask_EveryDay
	case utils.EveryMondayAt12AM:
		return UpdateCoverageTask_EveryWeek
	case utils.EveryMonthOn1stAt12AM:
		return UpdateCoverageTask_EveryMonth
	default:
		return func() {}
	}
}

func UpdateCoverageTask_EveryMinute() {
	fmt.Println("test")
}

func UpdateCoverageTask_EveryDay() {
	tasks, err := getTaskSchedules(utils.EveryDayAt12AM)
	if err != nil {
		log.Printf("Worker Error: %+v", err)
		return
	}

	if len(tasks) == 0 {
		jobCleanup(utils.EveryDayAt12AM)
		return
	}

	updateCoverageForTasks(tasks)

}

func UpdateCoverageTask_EveryWeek() {
	tasks, err := getTaskSchedules(utils.EveryMondayAt12AM)
	if err != nil {
		log.Printf("Worker Error: %+v", err)
		return
	}

	if len(tasks) == 0 {
		jobCleanup(utils.EveryMondayAt12AM)
		return
	}

	updateCoverageForTasks(tasks)
}

func UpdateCoverageTask_EveryMonth() {
	tasks, err := getTaskSchedules(utils.EveryMonthOn1stAt12AM)
	if err != nil {
		log.Printf("Worker Error: %+v", err)
		return
	}

	if len(tasks) == 0 {
		jobCleanup(utils.EveryMonthOn1stAt12AM)
		return
	}

	updateCoverageForTasks(tasks)
}

func updateCoverageForTasks(tasks []models.TaskSchedule) {
	for _, task := range tasks {
		directory, err := getDirectory(task.DirectoryID)
		if err != nil {
			log.Printf("Worker Error: %+v", err)
			return
		}

		err = utils.GenerateCoverageAndHTMLFiles(directory.Path)
		if err != nil {
			log.Printf("Worker Error: %+v", err)
			return
		}
	}
}

func jobCleanup(scheduleType utils.CronJobScheduleType) error {
	job, err := getJobsByScheduleType(scheduleType)
	if err != nil {
		return errors.WithStack(err)
	}

	if err := Clients.DB.Transaction(func(tx *gorm.DB) error {
		if err := deleteCronJob(job.ID); err != nil {
			return errors.WithStack(err)
		}

		Clients.Cron.RemoveCronJob(job.EntryID)

		return nil
	}); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func getTaskSchedules(scheduleType utils.CronJobScheduleType) ([]models.TaskSchedule, error) {
	var tasks []models.TaskSchedule

	result := Clients.DB.Joins("Directory").Where("schedule_type = ?", scheduleType).Find(&tasks)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return tasks, nil
}

func getDirectory(id uuid.UUID) (*models.Directory, error) {
	var directory *models.Directory

	result := Clients.DB.Model(&models.Directory{}).Where("id = ?", id).Find(&directory)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return directory, nil
}

func getJobsByScheduleType(scheduleType utils.CronJobScheduleType) (*models.CronJob, error) {
	var cronJob *models.CronJob

	result := Clients.DB.Model(&models.CronJob{}).Where("schedule_type = ?", scheduleType).Find(&cronJob)
	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return cronJob, nil
}

func deleteCronJob(id uuid.UUID) error {
	var cronJob *models.CronJob
	result := Clients.DB.Model(&models.CronJob{}).Where("id = ?", id).Delete(&cronJob)

	if result.Error != nil {
		return errors.WithStack(result.Error)
	}

	return nil
}
