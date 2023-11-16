package ignored

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"strconv"

	"github.com/gin-gonic/gin"
)

type CreateIgnoredRequest struct {
	Name       string `json:"name"`
	IgnoreType string `json:"ignoreType"`
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

	err = models.CreateIgnored(
		ctx,
		clients.DB,
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
