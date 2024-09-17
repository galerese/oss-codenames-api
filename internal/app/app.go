package app

import (
	"galere.se/oss-codenames-api/configs"
)

func Run(cfg *configs.Config) {
	// l := logging.New(cfg.LogLevel, cfg.LogFormat)

	// // Setup router
	// router := gin.Default()

	// gameRepository := database.NewDatabase()

	// // Setup routes
	// gameController := gamec.NewController(gamed.NewService(gameRepository), l)
	// gameController.SetupRoutes(router)

	// // Start server
	// router.Run(fmt.Sprintf(":%d", cfg.Port))
}
