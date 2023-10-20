package tasks

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
	Message string                   `json:"message"`
	Tasks   []map[string]interface{} `json:"tasks"`
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

	var tasks []map[string]interface{}
	for index, directory := range directories {
		taskSchedule, err := models.CreateTaskSchedule(ctx, clients.DB, directory, message.ScheduleType)
		if err != nil {
			return nil, api.InternalServerError(err)
		}

		task := map[string]interface{}{
			"ID":               taskSchedule.ID,
			"DirectoryID":      taskSchedule.DirectoryID,
			"ScheduleType":     taskSchedule.ScheduleType,
			"id":               index + 1,
			"coverageName":     directory.CoverageName,
			"scheduleTypeName": getScheduleTypeName(taskSchedule.ScheduleType),
		}

		tasks = append(tasks, task)
	}

	return &CreateTasksResponse{
		Message: "Good!",
		Tasks:   tasks,
	}, nil
}

func getScheduleTypeName(scheduleType utils.CronJobScheduleType) string {
	switch scheduleType {
	case 1:
		return "Daily"
	case 2:
		return "Weekly"
	case 3:
		return "Monthly"
	default:
		return ""
	}
}
