package settings

import (
	"golens-api/api"
	"golens-api/clients"
	client_funcs "golens-api/clients"
	"golens-api/models"
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CreateTaskRequest struct {
	DirectoryID  uuid.UUID                 `json:"directoryId" validate:"required"`
	ScheduleType utils.CronJobScheduleType `json:"scheduleType"`
}

type CreateTaskResponse struct {
	Message      string               `json:"message"`
	Task         *models.TaskSchedule `json:"task"`
	CoverageName string               `json:"coverageName"`
}

func CreateTask(
	ctx *gin.Context,
	message *CreateTaskRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	// TODO: there is a bug with the validate tag when you work with int
	// Look into if this is package related, maybe open a PR.
	if message.ScheduleType == 0 {
		return nil, api.InternalServerError(api.BadRequest)
	}

	err := handleJobCreation(ctx, clients, message.ScheduleType)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	directory, _, err := models.GetDirectory(ctx, clients.DB, message.DirectoryID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	taskSchedule, err := models.CreateTaskSchedule(ctx, clients.DB, *directory, message.ScheduleType)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	return &CreateTaskResponse{
		Message:      "Good!",
		Task:         taskSchedule,
		CoverageName: directory.CoverageName,
	}, nil
}

func handleJobCreation(ctx *gin.Context, clients *clients.GlobalClients, scheduleType utils.CronJobScheduleType) error {
	exists := models.CronJobExists(ctx, clients.DB, scheduleType)

	if !exists {
		err := clients.DB.Transaction(func(tx *gorm.DB) error {
			entryID, err := clients.Cron.CreateCronJob(scheduleType, client_funcs.GetUpdateTaskFunc(scheduleType))
			if err != nil {
				return errors.WithStack(err)
			}

			_, err = models.CreateCronJob(ctx, tx, scheduleType, entryID)
			if err != nil {
				return errors.WithStack(err)
			}

			return nil
		})
		if err != nil {
			return errors.WithStack(err)
		}
	}
	return nil
}
