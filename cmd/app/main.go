package main

import (
	"commentsService/internal/app"
	"commentsService/internal/config"
	"errors"
	"log"

	"github.com/joho/godotenv"
)

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
