package utils

import "github.com/gin-gonic/gin"

type Header struct {
	Auth      string
	Origin    string
	IpAddress string
}

const (
	Authorization = "Authorization"
	Origin        = "Origin"
)

func GetAPIHeaders(ctx *gin.Context) Header {
	return Header{
		Auth:      ctx.GetHeader(Authorization),
		Origin:    ctx.GetHeader(Origin),
		IpAddress: ctx.ClientIP(),
	}
}
