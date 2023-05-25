package clients

import (
	"golens-api/models"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"github.com/robfig/cron/v3"
)

type Task struct {
	ID       cron.EntryID
	Schedule string
	Handler  func()
}

type Cron struct {
	CronScheduler *cron.Cron
	Tasks         []Task
}

func InitializeCron() (*Cron, error) {
	tasks, err := findRunningTasks()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	cron := &Cron{
		CronScheduler: cron.New(),
		Tasks:         tasks,
	}

	cron.CronScheduler.Start()
	return cron, nil
}

func findRunningTasks() ([]Task, error) {
	var tasks []Task

	tasks, err := models.GetTaskSchedules(&gin.Context{}, Clients.DB)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return tasks, nil
}
