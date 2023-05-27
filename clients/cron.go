package clients

import (
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

func (c *Cron) ApplyRunningJobs() error {
	jobs, err := getJobs()
	if err != nil {
		return errors.WithStack(err)
	}

	for _, job := range jobs {
		newEntryID, err := c.CronScheduler.AddFunc(job.Schedule, GetUpdateTaskFunc(job.ScheduleType))
		if err != nil {
			return errors.WithStack(err)
		}

		job.EntryID = newEntryID

		err = updateJob(job)
		if err != nil {
			return errors.WithStack(err)
		}
	}

	return nil
}

func (c *Cron) GetEntries() []cron.Entry {
	crons := c.CronScheduler.Entries()
	return crons
}

func getJobs() ([]models.CronJob, error) {
	var jobs []models.CronJob

	result := Clients.DB.Model(&models.CronJob{}).Find(&jobs)

	if result.Error != nil {
		return nil, errors.WithStack(result.Error)
	}

	return jobs, nil
}

func updateJob(cronJob models.CronJob) error {
	result := Clients.DB.Model(&models.CronJob{}).Where("id = ?", cronJob.ID).Updates(&cronJob)
	if result.Error != nil {
		return errors.WithStack(result.Error)
	}
	return nil
}
