package settings

import (
	"fmt"
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TestRequest struct {
	ID uuid.UUID `json:"id"`
}

type TestResponse struct {
	Message string `json:"message"`
}

func Test(
	ctx *gin.Context,
	message *TestRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	// cronScheduler = cron.New()
	// cronScheduler.Start()
	// tasks = []Task{}
	// newTask := Task{
	// 	ID:       1,
	// 	Schedule: "* * * * *",
	// 	Handler: func() {
	// 		fmt.Println("hello")
	// 	},
	// }

	// cronScheduler.AddFunc(newTask.Schedule, newTask.Handler)
	// tasks = append(tasks, newTask)
	// fmt.Println(tasks)

	task, _ := models.GetTaskScheduleByDirectoryID(ctx, clients.DB, message.ID)
	fmt.Println(task.ID)

	return &TestResponse{
		Message: "Good!",
	}, nil
}
