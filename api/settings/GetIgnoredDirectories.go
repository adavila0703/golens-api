package settings

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type GetIgnoredDirectoriesRequest struct {
}

type GetIgnoredDirectoriesResponse struct {
	Directories []string `json:"directories"`
	Message     string   `json:"message"`
}

func GetIgnoredDirectories(
	ctx *gin.Context,
	message *GetIgnoredDirectoriesRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	directories := models.GetIgnoredDirectories(ctx, clients.DB)

	return &GetIgnoredDirectoriesResponse{
		Directories: directories,
		Message:     "Good!",
	}, nil
}
