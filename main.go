package main

import (
	"carefinder/database"
	"carefinder/routes"
	"log"

	"github.com/joho/godotenv"
)

func init() {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}
}

func main() {
	// Initialize Database
	database.Connect()

	// Setup Router
	r := routes.SetupRouter()

	// Start Server
	r.Run(":8080")
}
