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
	"gorm.io/gorm"
)

type CreateDirectoriesRequest struct {
	RootPath string `json:"rootPath" validate:"required"`
}

type CreateDirectoriesResponse struct {
	Message string `json:"message"`
}

func CreateDirectories(
	ctx *gin.Context,
	message *CreateDirectoriesRequest,
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

	for _, path := range paths {
		isGoDir, err := utils.IsGoDirectory(path)
		if err != nil {
			return nil, &api.Error{
				Err:    err,
				Status: http.StatusInternalServerError,
			}
		}

		if isGoDir {
			err := clients.DB.Transaction(func(tx *gorm.DB) error {
				err := models.CreateDirectory(ctx, tx, path)
				if err != nil {
					return errors.WithStack(err)
				}

				err = utils.GenerateCoverageAndHTMLFiles(path)
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
		}
	}

	return &CreateDirectoriesResponse{
		Message: "Good!",
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
