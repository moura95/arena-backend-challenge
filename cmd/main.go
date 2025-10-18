package main

import (
	"log"

	"arena-backend-challenge/config"
	server "arena-backend-challenge/internal"
)

// @title IP Location API
// @version 1.0.0
// @description REST API that resolves IP addresses to geographic locations using IP2Location dataset
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /
// @schemes http https

func main() {
	if err := config.LoadEnvFile(".env"); err != nil {
		log.Println("Warning: .env file not found, using defaults")
	}

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	srv, err := server.NewServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := srv.Start(); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
