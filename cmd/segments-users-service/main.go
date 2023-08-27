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

func main() {
	app.Run()
}