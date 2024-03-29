package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type UpdateDirectoryRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type UpdateDirectoryResponse struct {
	Message   string         `json:"message"`
	Directory map[string]any `json:"directory"`
}

func UpdateDirectory(
	ctx *gin.Context,
	message *UpdateDirectoryRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	directory, found, err := models.GetDirectory(ctx, clients.DB, message.ID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if !found {
		return nil, nil
	}

	err = clients.Cov.GenerateCoverageAndHTMLFiles(directory.Path)
	if err != nil {
		return nil, api.InternalServerError(err)
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

	return &UpdateDirectoryResponse{
		Message:   "Good!",
		Directory: directoryMap,
	}, nil
}
