package ignored

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type DeleteIgnoredRequest struct {
	ID uuid.UUID `json:"id"`
}

type DeleteIgnoredResponse struct {
	Message string `json:"message"`
}

func DeleteIgnored(
	ctx *gin.Context,
	message *DeleteIgnoredRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	err := models.DeleteIgnored(ctx, clients.DB, message.ID)
	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &DeleteIgnoredResponse{
		Message: "Good!",
	}, nil
}
