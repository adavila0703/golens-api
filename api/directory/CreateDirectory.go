package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CreateDirectoryRequest struct {
	Path string `json:"path" validate:"required"`
}

type CreateDirectoryResponse struct {
	Message string `json:"message"`
}

func CreateDirectory(
	ctx *gin.Context,
	message *CreateDirectoryRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	err := clients.DB.Transaction(func(tx *gorm.DB) error {
		err := models.CreateDirectory(ctx, tx, message.Path)
		if err != nil {
			return errors.WithStack(err)
		}

		err = utils.GenerateCoverageAndHTMLFiles(message.Path)
		if err != nil {
			return errors.WithStack(err)
		}

		return nil
	})
	if err != nil {
		return nil, &api.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	return &CreateDirectoryResponse{
		Message: "Good!",
	}, nil
}
