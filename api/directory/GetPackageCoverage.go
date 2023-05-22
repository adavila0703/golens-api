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

	packageMap, err := utils.GetPackageCoveragePercentage(directory.CoverageName)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	packageCoverage := createResponse(packageMap)

	return &GetPackageCoverageResponse{
		Message:         "Good!",
		PackageCoverage: packageCoverage,
	}, nil
}

func createResponse(
	packageMap map[string]map[string]int,
) []map[string]any {
	var packageCoverage []map[string]any
	for packageName, packageValue := range packageMap {
		percentage := utils.GetCoveragePercentageNumber(packageValue["totalStatements"], packageValue["coveredStatements"])
		coverage := map[string]any{
			"packageName": packageName,
			"coverage":    percentage,
		}

		packageCoverage = append(packageCoverage, coverage)
	}

	return packageCoverage
}
