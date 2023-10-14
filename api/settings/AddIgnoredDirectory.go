package settings

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type AddIgnoredDirectoryRequest struct {
	Directory string `json:"directory"`
}

type AddIgnoredDirectoryResponse struct {
	Message string `json:"message"`
}

func AddIgnoredDirectory(
	ctx *gin.Context,
	message *AddIgnoredDirectoryRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	err := models.AddIgnoredDirectory(ctx, clients.DB, message.Directory)

	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &AddIgnoredDirectoryResponse{
		Message: "Good!",
	}, nil
}
