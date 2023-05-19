package health

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.Engine, group string) {
	router.
		Group(group).
		GET("health", api.Handler(HealthCheck))
}
