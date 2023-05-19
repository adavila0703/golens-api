package middleware

import (
	"errors"
	"golens-api/clients"
	"golens-api/config"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	redis "github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

func Limiter() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		redisClient := clients.Clients.Redis
		ipAddress := ctx.ClientIP()

		value, err := getIpKeyValueRedis(ctx, redisClient, ipAddress)
		if err != nil {
			logrus.WithFields(logrus.Fields{"ipaddress": ipAddress}).Error(err.Error())
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		if value == config.Cfg.MaxRequests {
			ctx.Abort()
		} else {
			duration := time.Duration(5) * time.Second
			redisClient.Set(ctx, ipAddress, value+1, duration)
		}

		ctx.Next()
	}
}

func getIpKeyValueRedis(ctx *gin.Context, redisClient *redis.Client, ipAddress string) (int64, error) {
	value := redisClient.Get(ctx, ipAddress).Val()
	if value == "" {
		return 0, nil
	}

	parsedValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return 0, errors.New("error parsing redis value")
	}

	return parsedValue, nil
}
