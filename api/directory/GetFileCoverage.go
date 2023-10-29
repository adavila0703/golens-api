package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetFileCoverageRequest struct {
	RepoID      uuid.UUID `json:"repoId" validate:"required"`
	PackageName string    `json:"packageName" validate:"required"`
}

type GetFileCoverageResponse struct {
	Message      string           `json:"message"`
	FileCoverage []map[string]any `json:"fileCoverage"`
}

func GetFileCoverage(
	ctx *gin.Context,
	message *GetFileCoverageRequest,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	directory, found, err := models.GetDirectory(ctx, clients.DB, message.RepoID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if !found {
		return &GetFileCoverageResponse{
			Message: "Directory not found",
		}, nil
	}

	fileCoverage, err := clients.Utils.GetFileCoveragePercentage(directory.CoverageName)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	return &GetFileCoverageResponse{
		Message:      "Good!",
		FileCoverage: fileCoverage[message.PackageName],
	}, nil
}
