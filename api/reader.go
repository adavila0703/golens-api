package api

import (
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	jsoniter "github.com/json-iterator/go"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

// reads a message of type T
func ReadRequest[T any](ctx *gin.Context) (*T, error) {
	var message *T
	bytes, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		log.
			WithFields(log.Fields{"stack": "ReadRequest_ReadAll"}).
			Error("error reading request")
		ctx.AbortWithStatus(http.StatusBadRequest)
		return nil, errors.WithStack(err)
	}

	err = jsoniter.Unmarshal(bytes, &message)
	if err != nil {
		log.
			WithFields(log.Fields{"stack": "ReadRequest_Unmarshal"}).
			Error("error reading request", err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return nil, errors.WithStack(err)
	}

	// validate fields
	validate := validator.New()
	err = validate.StructCtx(ctx, message)
	if err != nil {
		log.
			WithFields(log.Fields{"stack": "ReadRequest"}).
			Error(err)
		ctx.AbortWithStatus(http.StatusBadRequest)
		return nil, errors.WithStack(err)
	}

	return message, nil
}
