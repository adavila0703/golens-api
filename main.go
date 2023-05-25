package main

import (
	"golens-api/api/directory"
	"golens-api/api/health"
	"golens-api/api/settings"
	"golens-api/clients"
	"golens-api/config"
	"golens-api/middleware"
	"golens-api/models"
	"log"

	"github.com/gin-gonic/gin"

	"github.com/gin-contrib/cors"
	"github.com/kelseyhightower/envconfig"
)

func main() {
	// process environment variables
	err := envconfig.Process("", &config.Cfg)
	if err != nil {
		log.Fatalf("Error processing env vars: %s", err)
	}

	// initialize postgres client
	postgres, err := clients.NewPostgresClient(config.Cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Postgres error: %s", err)
	}

	// initialize redis client
	// redisClient := clients.NewRedisClient(&redis.Options{
	// 	Addr:     config.Cfg.RediscloudURL,
	// 	Username: config.Cfg.RedisUsername,
	// 	Password: config.Cfg.RedisPassword,
	// })

	// redisPing := redisClient.Ping(context.Background())
	// if redisPing.Err() != nil {
	// 	log.Fatalf("Redis error: %s", redisPing.Err())
	// }

	cron := clients.InitializeCron()

	// initialize global clients
	clients.Clients = clients.NewGlobalClients(
		postgres,
		nil,
		cron,
	)

	// migrate db models
	models.MigrateModels(postgres)

	// set up router
	router := gin.Default()
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{config.Cfg.AllowOrigin},
		AllowMethods:     []string{"PUT", "PATCH", "GET"},
		AllowHeaders:     []string{"Origin", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// api key authentication
	router.Use(middleware.Auth())

	// rate limiter middleware
	// router.Use(middleware.Limiter())

	// health route
	health.SubRoutes(router, "")

	// api sub routes
	apiRouter := router.Group("api")
	directory.SubRoutes(apiRouter, "directory")
	settings.SubRoutes(apiRouter, "settings")

	if config.Cfg.HostPort != "" {
		err = router.Run(config.Cfg.HostPort)
	} else {
		err = router.Run()
	}

	if err != nil {
		log.Fatalf("Gin error: %s", err)
	}
}
