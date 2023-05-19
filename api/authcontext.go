package api

import (
	"golens-api/clients"
	"golens-api/utils"

	"github.com/gin-gonic/gin"
)

type AuthContext struct {
	Username  string
	IpAddress string
	Headers   utils.Header
}

// gets the auth context of the incoming request
func GetAuthContext(ctx *gin.Context, clients *clients.GlobalClients) *AuthContext {
	// db := clients.DB
	headers := utils.GetAPIHeaders(ctx)

	authContext := &AuthContext{
		Username:  "",
		Headers:   headers,
		IpAddress: headers.IpAddress,
	}

	// var user *models.User

	// result := db.WithContext(ctx).
	// 	Model(&models.User{}).
	// 	Select("device_uuid, username").
	// 	Where("device_uuid = ?", headers.DeviceUUID).
	// 	First(&user)
	// if result.RowsAffected > 0 {
	// 	authContext.Username = user.Username
	// } else {
	// 	authContext.Username = utils.RandoName()
	// }

	return authContext
}
