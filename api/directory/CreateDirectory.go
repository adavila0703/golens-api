package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type CreateDirectoryRequest struct {
	Path string `json:"path" validate:"required"`
}

type CreateDirectoryResponse struct {
	Message   string         `json:"message"`
	Directory map[string]any `json:"directory"`
}

func CreateDirectory(
	ctx *gin.Context,
	message *CreateDirectoryRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	found, err := models.DirectoryExists(ctx, clients.DB, message.Path)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	// TODO: change this to return an error which the frontend can handle
	if found {
		return nil, nil
	}

	isGoDirectory, err := clients.Cov.IsGoDirectory(message.Path)
	if !isGoDirectory || err != nil {
		if err != nil {
			return nil, &api.Error{
				Err:    err,
				Status: http.StatusInternalServerError,
			}
		}

		if !isGoDirectory {
			return nil, &api.Error{
				Err:    errors.New("Is not a go directory"),
				Status: http.StatusBadRequest,
			}
		}
	}

	var directory *models.Directory
	err = clients.DB.Transaction(func(tx *gorm.DB) error {
		var err error
		directory, err = models.CreateDirectory(ctx, tx, message.Path)
		if err != nil {
			return errors.WithStack(err)
		}

		err = clients.Cov.GenerateCoverageAndHTMLFiles(message.Path)
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

	ignoredPackages := clients.Cov.GetIgnoredPackages(ctx, clients.DB, directory.CoverageName)

	totalLines, coveredLines, err := clients.Cov.GetCoveredLines(directory.CoverageName, ignoredPackages)
	if err != nil {
		return nil, &api.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	directoryMap := map[string]any{
		"id":           directory.ID.String(),
		"path":         directory.Path,
		"totalLines":   totalLines,
		"coveredLines": coveredLines,
		"coverageName": directory.CoverageName,
	}

	return &CreateDirectoryResponse{
		Message:   "Good!",
		Directory: directoryMap,
	}, nil
}
