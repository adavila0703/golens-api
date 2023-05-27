package settings

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type DeleteTaskRequest struct {
	TaskID       uuid.UUID                 `json:"taskID" validate:"required"`
	ScheduleType utils.CronJobScheduleType `json:"scheduleType"`
}

type DeleteTaskResponse struct {
	Message string `json:"message"`
}

func DeleteTask(
	ctx *gin.Context,
	message *DeleteTaskRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	if err := models.DeleteTaskSchedule(ctx, clients.DB, message.TaskID); err != nil {
		return nil, api.InternalServerError(err)
	}

	tasks, err := models.GetTaskSchedulesByScheduleType(ctx, clients.DB, message.ScheduleType)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if len(tasks) == 0 {
		job, err := models.GetCronJob(ctx, clients.DB, message.ScheduleType)
		if err != nil {
			return nil, api.InternalServerError(err)
		}

		if err := clients.DB.Transaction(func(tx *gorm.DB) error {
			if err := models.DeleteCronJob(ctx, tx, job.ID); err != nil {
				return errors.WithStack(err)
			}

			clients.Cron.RemoveCronJob(job.EntryID)

			return nil
		}); err != nil {
			return nil, api.InternalServerError(err)
		}
	}

	return &DeleteTaskResponse{
		Message: "Good!",
	}, nil
}
