package settings

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type CreateIgnoredDirectoryRequest struct {
	DirectoryName string `json:"directoryName"`
}

type CreateIgnoredDirectoryResponse struct {
	Message string `json:"message"`
}

func CreateIgnoredDirectory(
	ctx *gin.Context,
	message *CreateIgnoredDirectoryRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	err := models.CreateIgnoredDirectory(ctx, clients.DB, message.DirectoryName)

	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &CreateIgnoredDirectoryResponse{
		Message: "Good!",
	}, nil
}
