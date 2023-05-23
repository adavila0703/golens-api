package directory

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.RouterGroup, group string) {
	router.
		Group(group).
		GET("GetDirectories", api.Handler(GetDirectories)).
		POST("GetHtmlContents", api.Handler(GetHtmlContents)).
		POST("DeleteDirectory", api.Handler(DeleteDirectory)).
		POST("GetFileCoverage", api.Handler(GetFileCoverage)).
		POST("GetPackageCoverage", api.Handler(GetPackageCoverage)).
		POST("CreateDirectory", api.Handler(CreateDirectory)).
		POST("CreateDirectories", api.Handler(CreateDirectories))
}
