package main

import (
	"galere.se/oss-codenames-api/configs"
	"galere.se/oss-codenames-api/internal/app"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	logrus.Info("Starting OSS CodeNames API server")

	// Bootstrap environment config from .env.example file
	_ = godotenv.Load(".env")

	// Configuration
	cfg, err := configs.NewConfig()
	if err != nil {
		logrus.Fatalf("Config error: %s", err)
	}

	// Run
	app.Run(cfg)

}
