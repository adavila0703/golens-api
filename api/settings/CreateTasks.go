package settings

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"

	"github.com/gin-gonic/gin"
)

type CreateTasksRequest struct {
	ScheduleType utils.CronJobScheduleType `json:"scheduleType"`
}

type CreateTasksResponse struct {
	Message string `json:"message"`
}

func CreateTasks(
	ctx *gin.Context,
	message *CreateTasksRequest,
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

	directories, err := models.GetDirectories(ctx, clients.DB)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	for _, directory := range directories {
		_, err := models.CreateTaskSchedule(ctx, clients.DB, directory, message.ScheduleType)
		if err != nil {
			return nil, api.InternalServerError(err)
		}
	}

	return &CreateTasksResponse{
		Message: "Good!",
	}, nil
}
