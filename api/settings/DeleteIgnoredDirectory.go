package settings

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteIgnoredDirectoryRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteIgnoredDirectoryResponse struct {
	Message string `json:"message"`
}

func DeleteIgnoredDirectory(
	ctx *gin.Context,
	message *DeleteIgnoredDirectoryRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	err := models.DeleteIgnoredDirectory(ctx, clients.DB, message.ID)

	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &DeleteIgnoredDirectoryResponse{
		Message: "Good!",
	}, nil
}
