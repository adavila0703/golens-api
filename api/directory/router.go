package directory

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.RouterGroup, group string) {
	router.
		Group(group).
		GET("GetDirectories", api.Handler(GetDirectories)).
		POST("CreateDirectory", api.Handler(CreateDirectory)).
		POST("CreateDirectories", api.Handler(CreateDirectories))
}
