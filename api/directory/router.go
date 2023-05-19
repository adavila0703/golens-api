package directory

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.RouterGroup, group string) {
	router.
		Group(group).
		POST("CreateDirectory", api.Handler(CreateDirectory)).
		GET("GetDirectories", api.Handler(GetDirectories))
}
