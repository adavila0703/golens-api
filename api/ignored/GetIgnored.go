package ignored

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type GetIgnoredRequest struct {
}

type GetIgnoredResponse struct {
	Directories []models.Ignored `json:"directories"`
	Message     string           `json:"message"`
}

func GetIgnored(
	ctx *gin.Context,
	message *GetIgnoredRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	directories := models.GetIgnored(ctx, clients.DB, models.DirectoryType)

	return &GetIgnoredResponse{
		Directories: directories,
		Message:     "Good!",
	}, nil
}
