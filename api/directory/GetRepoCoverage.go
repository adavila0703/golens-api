package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type GetRepoCoverageRequest struct {
	ID uuid.UUID `json:"id" validate:"required"`
}

type GetRepoCoverageResponse struct {
	Message         string                      `json:"message"`
	PackageCoverage []map[string]any            `json:"packageCoverage"`
	FileCoverage    map[string][]map[string]any `json:"fileCoverage"`
}

func GetRepoCoverage(
	ctx *gin.Context,
	message *GetRepoCoverageRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {

	directory, found, err := models.GetDirectory(ctx, clients.DB, message.ID)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	if !found {
		return &GetRepoCoverageResponse{
			Message: "Directory not found",
		}, nil
	}

	fileCoverage, err := utils.GetFileCoveragePercentage(directory.CoverageName)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	packageMap, err := utils.GetPackageCoveragePercentage(directory.CoverageName)
	if err != nil {
		return nil, api.InternalServerError(err)
	}

	repoCoverage := getRepoCoverageResponse(packageMap)

	return &GetRepoCoverageResponse{
		Message:         "Good!",
		PackageCoverage: repoCoverage,
		FileCoverage:    fileCoverage,
	}, nil
}

func getRepoCoverageResponse(
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
