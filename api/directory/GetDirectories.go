package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type GetDirectoriesRequest struct {
}

type GetDirectoriesResponse struct {
	Message     string           `json:"message"`
	Directories []map[string]any `json:"directories"`
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

	var directoryMaps []map[string]any

	for _, directory := range directories {
		totalLines, coveredLines, err := utils.GetCoveredLines(directory.CoverageName)
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

		directoryMaps = append(directoryMaps, directoryMap)
	}

	return &GetDirectoriesResponse{
		Message:     "Good!",
		Directories: directoryMaps,
	}, nil
}
