package directory

import (
	"golens-api/api"
	"golens-api/clients"
	"golens-api/models"
	"golens-api/utils"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type GetRootDirectoryPathsRequest struct {
	RootPath string `json:"rootPath" validate:"required"`
}

type GetRootDirectoryPathsResponse struct {
	Message string   `json:"message"`
	Paths   []string `json:"paths"`
}

func GetRootDirectoryPaths(
	ctx *gin.Context,
	message *GetRootDirectoryPathsRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	paths, err := getDirPaths(message.RootPath)
	if err != nil {
		return nil, &api.Error{
			Err:    err,
			Status: http.StatusInternalServerError,
		}
	}

	ignoredDirectoriesMap := make(map[string]bool)
	ignoredDirectories := models.GetIgnoredDirectories(ctx, clients.DB)

	for _, directory := range ignoredDirectories {
		ignoredDirectoriesMap[directory] = true
	}

	var goPaths []string
	for _, path := range paths {
		isGoDir, err := utils.IsGoDirectory(path)
		if err != nil {
			return nil, &api.Error{
				Err:    err,
				Status: http.StatusInternalServerError,
			}
		}

		directoryName := utils.GetCoverageNameFromPath(path)
		_, ok := ignoredDirectoriesMap[directoryName]

		if isGoDir && !ok {
			goPaths = append(goPaths, path)
		}
	}

	return &GetRootDirectoryPathsResponse{
		Message: "Good!",
		Paths:   goPaths,
	}, nil
}

func getDirPaths(rootPath string) ([]string, error) {
	var paths []string

	dirEntries, err := os.ReadDir(rootPath)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	for _, dirEntry := range dirEntries {
		if dirEntry.IsDir() {
			dirPath := filepath.Join(rootPath, dirEntry.Name())

			paths = append(paths, dirPath)
		}
	}

	return paths, nil
}
