package ignored

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type CreateIgnoredRequest struct {
	DirectoryID uuid.UUID `json:"directoryId"`
	Name        string    `json:"name"`
	IgnoreType  string    `json:"ignoreType"`
}

type CreateIgnoredResponse struct {
	Message string `json:"message"`
}

func CreateIgnored(
	ctx *gin.Context,
	message *CreateIgnoredRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	ignoreType, err := strconv.ParseInt(message.IgnoreType, 10, 64)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	directory, found, err := models.GetDirectory(ctx, clients.DB, message.DirectoryID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if !found {
		return nil, nil
	}

	err = models.CreateIgnored(
		ctx,
		clients.DB,
		directory.CoverageName,
		message.Name,
		models.IgnoreType(ignoreType),
	)

	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &CreateIgnoredResponse{
		Message: "Good!",
	}, nil
}
