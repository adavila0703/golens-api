package ignored

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type CreateIgnoredRequest struct {
	DirectoryName string `json:"directoryName"`
}

type CreateIgnoredResponse struct {
	Message string `json:"message"`
}

func CreateIgnored(
	ctx *gin.Context,
	message *CreateIgnoredRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	err := models.CreateIgnored(ctx, clients.DB, message.DirectoryName, models.DirectoryType)

	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &CreateIgnoredResponse{
		Message: "Good!",
	}, nil
}
