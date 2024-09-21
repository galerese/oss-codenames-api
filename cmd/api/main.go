package main

import (
	"log"

	"galere.se/oss-codenames-api/configs"
	"galere.se/oss-codenames-api/internal/app"
	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting OSS CodeNames API server")

	// Bootstrap environment config from .env.dev file
	_ = godotenv.Load(".env")

	// Configuration
	cfg, err := configs.NewConfig()
	if err != nil {
		log.Fatalf("Config error: %s", err)
	}

	// Run
	err = app.Run(cfg)
	if err != nil {
		log.Fatalf("Error running app: %s", err)
	}
}
