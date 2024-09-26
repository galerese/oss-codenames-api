package app

import (
	"context"
	"fmt"
	"time"

	"galere.se/oss-codenames-api/configs"
	gamec "galere.se/oss-codenames-api/internal/controller/game"
	sessionc "galere.se/oss-codenames-api/internal/controller/session"
	gamed "galere.se/oss-codenames-api/internal/domain/game"
	sessiond "galere.se/oss-codenames-api/internal/domain/session"
	"galere.se/oss-codenames-api/internal/repository/database"
	"galere.se/oss-codenames-api/pkg/http_controller"
	"galere.se/oss-codenames-api/pkg/logging"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func Run(cfg *configs.Config) error {
	l := logging.New(cfg.LogLevel, cfg.LogFormat)

	// Setup router
	router := gin.Default()
	http_controller.Bootstrap(router, l, cfg.AppName)

	// Initialize DB
	l.Info("Connecting to MongoDB...")
	mongoClient, err := mongo.Connect(options.Client().ApplyURI(cfg.DatabaseURL).SetTimeout(time.Second * 5))
	if err != nil {
		return errors.Wrap(err, "Could not initialize mongodb!")
	}

	err = mongoClient.Ping(context.Background(), nil)
	if err != nil {
		return errors.Wrap(err, "Could not connect to mongodb!")
	}

	gameRepository := database.NewDatabase(mongoClient.Database(cfg.DatabaseName), l)
	err = gameRepository.SetupLockExpiration(context.Background())
	if err != nil {
		return errors.Wrap(err, "Could not setup lock expiration for locking mechanics")
	}

	// Setup services
	l.Info("Initializing services...")
	sessionService := sessiond.NewService(gameRepository, &l)
	gameService := gamed.NewService(gameRepository, &l)

	// Setup controllers
	l.Info("Initializing controllers...")
	gameController := gamec.NewController(gameService, sessionService, gameRepository, &l)
	sessionController := sessionc.NewController(sessionService, &l)

	// Setup routes
	l.Info("Setting up routes...")
	rg := router.Group("/v1")
	gameController.SetupRoutes(rg)
	sessionController.SetupRoutes(rg)

	// Start server
	l.Info("Starting server...")
	router.Run(fmt.Sprintf(":%s", cfg.HttpPort))
	l.Info("Server started on port [%d] :)", cfg.HttpPort)

	return nil
}
