package settings

import (
	"fmt"
	"golens-api/api"
	"golens-api/clients"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/robfig/cron/v3"
)

type TestRequest struct {
	ID     uuid.UUID    `json:"id"`
	Num    int          `json:"num"`
	CronID cron.EntryID `json:"cron"`
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

	entries := clients.Cron.GetEntries()

	for _, entry := range entries {
		fmt.Println(entry.ID)
	}

	return &TestResponse{
		Message: "Good!",
	}, nil
}
