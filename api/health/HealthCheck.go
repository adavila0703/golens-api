package health

import (
	"golens-api/api"
	"golens-api/clients"

	"github.com/gin-gonic/gin"
)

type HealthCheckRequest struct {
}

type HealthCheckResponse struct {
	Message string `json:"message"`
}

func HealthCheck(
	ctx *gin.Context,
	message *HealthCheckRequest,
	authContext *api.AuthContext,
	clients *clients.GlobalClients,
) (interface{}, *api.Error) {
	return &HealthCheckResponse{Message: "Good!"}, nil
}
