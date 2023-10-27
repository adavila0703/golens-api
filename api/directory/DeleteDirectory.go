package directory

import (
	"fmt"
	"golens-api/api"
	"golens-api/api/tasks"
	"golens-api/clients"
	"golens-api/models"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteDirectoryRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type DeleteDirectoryResponse struct {
	Message string `json:"message"`
}

func DeleteDirectory(
	ctx *gin.Context,
	message *DeleteDirectoryRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	directory, found, err := models.GetDirectory(ctx, clients.DB, message.ID)
	if !found {
		return nil, api.InternalServerError(err)
	}

	err = models.DeleteDirectory(ctx, clients.DB, message.ID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	workingDir, err := os.Getwd()
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	coverageProfile := fmt.Sprintf("%s/data/coverage/%s.out", workingDir, directory.CoverageName)
	htmlFile := fmt.Sprintf("%s/data/html/%s.html", workingDir, directory.CoverageName)

	err = os.Remove(coverageProfile)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	err = os.Remove(htmlFile)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	task, err := models.GetTaskScheduleByDirectoryID(ctx, clients.DB, message.ID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	deleteTaskRequest := &tasks.DeleteTaskRequest{
		TaskID:       task.ID,
		ScheduleType: task.ScheduleType,
	}

	_, deleteTaskErr := tasks.DeleteTask(ctx, deleteTaskRequest, clients)
	if deleteTaskErr != nil {
		return nil, api.InternalServerError(err)
	}

	return &DeleteDirectoryResponse{
		Message: "Good!",
	}, nil
}
