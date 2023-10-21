package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetPackageCoverageRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type GetPackageCoverageResponse struct {
	Message         string           `json:"message"`
	PackageCoverage []map[string]any `json:"packageCoverage"`
}

func GetPackageCoverage(
	ctx *gin.Context,
	message *GetPackageCoverageRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	directory, found, err := models.GetDirectory(ctx, clients.DB, message.ID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if !found {
		return &GetPackageCoverageResponse{
			Message: "Directory not found",
		}, nil
	}

	coveredLinesByPackage, err := utils.GetCoveredLinesByPackage(directory.CoverageName)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	var packageCoverage []map[string]any
	for packageName, lines := range coveredLinesByPackage {
		coverage := map[string]any{
			"packageName":  packageName,
			"coveredLines": lines["coveredLines"],
			"totalLines":   lines["totalLines"],
		}

		packageCoverage = append(packageCoverage, coverage)
	}

	return &GetPackageCoverageResponse{
		Message:         "Good!",
		PackageCoverage: packageCoverage,
	}, nil
}
