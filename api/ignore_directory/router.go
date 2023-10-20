package ignore_directory

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.RouterGroup, group string) {
	router.
		Group(group).
		GET("GetIgnoredDirectories", api.Handler(GetIgnoredDirectories)).
		POST("CreateIgnoredDirectory", api.Handler(CreateIgnoredDirectory)).
		POST("DeleteIgnoredDirectory", api.Handler(DeleteIgnoredDirectory))
}
