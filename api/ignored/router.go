package ignored

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.RouterGroup, group string) {
	router.
		Group(group).
		GET("GetIgnored", api.Handler(GetIgnored)).
		POST("CreateIgnored", api.Handler(CreateIgnored)).
		POST("DeleteIgnored", api.Handler(DeleteIgnored))
}
