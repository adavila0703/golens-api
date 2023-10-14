package settings

import (
	"golens-api/api"

	"github.com/gin-gonic/gin"
)

func SubRoutes(router *gin.RouterGroup, group string) {
	router.
		Group(group).
		GET("GetTasks", api.Handler(GetTasks)).
		GET("GetIgnoredDirectories", api.Handler(GetIgnoredDirectories)).
		POST("AddIgnoredDirectory", api.Handler(AddIgnoredDirectory)).
		POST("DeleteTask", api.Handler(DeleteTask)).
		POST("CreateTask", api.Handler(CreateTask)).
		POST("DeleteTasks", api.Handler(DeleteTasks)).
		POST("CreateTasks", api.Handler(CreateTasks)).
		POST("Test", api.Handler(Test))
}
