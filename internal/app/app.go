package app

import (
	"fmt"

	"galere.se/oss-codenames-api/configs"
	gamec "galere.se/oss-codenames-api/internal/controller/game"
	sessionc "galere.se/oss-codenames-api/internal/controller/session"
	gamed "galere.se/oss-codenames-api/internal/domain/game"
	sessiond "galere.se/oss-codenames-api/internal/domain/session"
	"galere.se/oss-codenames-api/internal/repository/database"
	"galere.se/oss-codenames-api/pkg/http_controller"
	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-gonic/gin"
)

func Run(cfg *configs.Config) {
	l := logging.New(cfg.LogLevel, cfg.LogFormat)

	// Setup router
	router := gin.Default()
	http_controller.Bootstrap(router, l, cfg.AppName)

	gameRepository := database.NewDatabase()

	// Setup services
	sessionService := sessiond.NewService(gameRepository, &l)
	gameService := gamed.NewService(gameRepository, &l)

	// Setup controllers
	gameController := gamec.NewController(gameService, sessionService, &l)
	sessionController := sessionc.NewController(sessionService, &l)

	// Setup routes
	rg := router.Group("/v1")
	gameController.SetupRoutes(rg)
	sessionController.SetupRoutes(rg)

	// Start server
	router.Run(fmt.Sprintf(":%s", cfg.HttpPort))
	l.Info("Server started on port [%d] :)", cfg.HttpPort)
}
