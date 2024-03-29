package main

import (
	"golens-api/api/directory"
	"golens-api/api/health"
	"golens-api/api/ignored"
	"golens-api/api/tasks"
	"golens-api/clients"
	"golens-api/config"
	"golens-api/coverage"
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

	cron, err := clients.InitializeCron()
	if err != nil {
		log.Fatalf("Cron error: %s", err)
	}

	utilsClient := coverage.NewCoverage()

	// initialize global clients
	clients.Clients = clients.NewGlobalClients(
		postgres,
		cron,
		utilsClient,
	)

	// migrate db models
	models.MigrateModels(postgres)

	// find running tasks
	err = clients.Clients.Cron.ApplyRunningJobs()
	if err != nil {
		log.Fatalf("Cron error: %s", err)
	}

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

	// health route
	health.SubRoutes(router, "")

	// api sub routes
	apiRouter := router.Group("api")
	directory.SubRoutes(apiRouter, "directory")
	tasks.SubRoutes(apiRouter, "tasks")
	ignored.SubRoutes(apiRouter, "ignored")

	if config.Cfg.HostPort != "" {
		err = router.Run(config.Cfg.HostPort)
	} else {
		err = router.Run()
	}

	if err != nil {
		log.Fatalf("Gin error: %s", err)
	}
}
