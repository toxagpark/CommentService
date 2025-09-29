package main

import (
	"commentsService/internal/app"
	"commentsService/internal/config"
	"errors"
	"log"

	"github.com/joho/godotenv"
)

// @title Comment Service API
// @version 1.0
// @description HTML-based comment management API
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found. Using system environment variables")
	}

	cfg, err := config.NewConfig()
	if errors.Is(err, config.ErrConfig) {
		log.Println(err)
		return
	}
	app.App(cfg)
}
