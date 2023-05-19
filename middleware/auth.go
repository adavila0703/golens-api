package middleware

import (
	"golens-api/config"
	"golens-api/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

// middleware which require user authentication
func Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// auth middleware

		headers := utils.GetAPIHeaders(ctx)

		if headers.Origin != config.Cfg.AllowOrigin {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		if headers.Auth != config.Cfg.AllowApiKey {
			ctx.AbortWithStatus(http.StatusUnauthorized)
		}

		ctx.Next()
	}
}
