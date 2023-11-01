package tasks

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type GetTasksRequest struct {
}

type GetTasksResponse struct {
	Message string                `json:"message"`
	Tasks   []models.TaskSchedule `json:"tasks"`
}

func GetTasks(
	ctx *gin.Context,
	message *GetTasksRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	tasks, err := models.GetTaskSchedules(ctx, clients.DB)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	return &GetTasksResponse{
		Message: "Good!",
		Tasks:   tasks,
	}, nil
}
