package clients

import (
	"fmt"
	"golens-api/models"
	"golens-api/utils"

	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type Cron struct {
	CronScheduler *cron.Cron
}

func InitializeCron() (*Cron, error) {
	cron := &Cron{
		CronScheduler: cron.New(),
	}

	cron.CronScheduler.Start()

	return cron, nil
}

func (c *Cron) CreateCronJob(schedule utils.CronJobScheduleType, handler func()) (cron.EntryID, error) {
	cronSchedule := utils.GetCronSchedule(schedule)

	id, err := c.CronScheduler.AddFunc(cronSchedule, handler)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return id, nil
}

func (c *Cron) RemoveCronJob(id cron.EntryID) {
	c.CronScheduler.Remove(id)
}

// func (c *Cron) ApplyRunningTasks() error {
// 	taskSchedules, err := getTaskSchedules()
// 	if err != nil {
// 		return errors.WithStack(err)
// 	}

// 	return nil
// }

func printHello() {
	fmt.Println("hello")
}

func getTaskSchedules() ([]models.TaskSchedule, error) {
	var tasks []models.TaskSchedule

	result := Clients.DB.Model(&models.TaskSchedule{}).Find(&tasks)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return tasks, nil
}
