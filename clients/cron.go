package clients

import "github.com/robfig/cron/v3"

type Task struct {
	ID       cron.EntryID
	Schedule string
	Handler  func()
}

type Cron struct {
	CronScheduler *cron.Cron
	Tasks         []Task
}

func InitializeCron() *Cron {
	cron := &Cron{
		CronScheduler: cron.New(),
		Tasks:         []Task{},
	}

	cron.CronScheduler.Start()
	return cron
}
