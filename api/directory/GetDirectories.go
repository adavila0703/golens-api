package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
)

type GetDirectoriesRequest struct {
}

type GetDirectoriesResponse struct {
	Message     string             `json:"message"`
	Directories []models.Directory `json:"directories"`
}

func GetDirectories(
	ctx *gin.Context,
	message *GetDirectoriesRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	directories, err := models.GetDirectories(ctx, clients.DB)
	if err != nil {
		return nil, &api.Error{
			Err: err,
		}
	}

	return &GetDirectoriesResponse{
		Message:     "Good!",
		Directories: directories,
	}, nil
}
