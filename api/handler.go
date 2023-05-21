package api

import (
	"fmt"
	"golens-api/clients"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Err    error
	Status int
}

func InternalServerError(err error) *Error {
	return &Error{
		Err:    err,
		Status: http.StatusInternalServerError,
	}
}

// runs endpoint function
func Handler[T any](handleFunc func(*gin.Context, *T, *AuthContext, *clients.GlobalClients) (interface{}, *Error)) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		var err error
		var clients = clients.Clients

		authContext := GetAuthContext(ctx, clients)

		// read the incoming request message and validate fields
		var message *T
		if ctx.Request.Method != http.MethodGet {
			message, err = ReadRequest[T](ctx, authContext)
			if err != nil {
				log.
					WithFields(log.Fields{"user_id": authContext.Username, "stack": "Handler"}).
					Error("read request error")
				return
			}
		}

		// run the handle func
		payload, handlerError := handleFunc(ctx, message, authContext, clients)
		if handlerError != nil {
			if authContext != nil {
				log.
					WithFields(log.Fields{"user_id": authContext.Username, "stack": "Handler"}).
					Error("response error")
			}

			if handlerError.Err != nil {
				fmt.Printf("\n%+v\n", handlerError.Err)
			}

			if handlerError.Status != 0 {
				ctx.AbortWithStatus(handlerError.Status)
				return
			}

			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		ctx.JSON(http.StatusOK, payload)
	}
}
