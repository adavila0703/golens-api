package tasks

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type DeleteTasksRequest struct {
}

type DeleteTasksResponse struct {
	Message string `json:"message"`
}

func DeleteTasks(
	ctx *gin.Context,
	message *DeleteTasksRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	tasks, err := models.GetTaskSchedules(ctx, clients.DB)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	for _, task := range tasks {
		deleteTaskRequest := &DeleteTaskRequest{
			TaskID:       task.ID,
			ScheduleType: task.ScheduleType,
		}

		_, err := DeleteTask(ctx, deleteTaskRequest, clients)
		if err != nil {
			return nil, api.InternalServerError(err.Err)
		}
	}
	return &DeleteTasksResponse{
		Message: "Good!",
	}, nil
}
