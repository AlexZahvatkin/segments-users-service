package main

import (
	"log"

	"github.com/AlexZahvatkin/segments-users-service/internal/app"
	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

// @title Segments Users Service
// @version 1.0
// @description API Server segments management.
// @host localhost:8080
// @BasePath /v1
// @accept json
// @produce json
// @schemes http

func main() {
	app.Run()
}