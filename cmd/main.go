package main

import (
	"log"

	"arena-backend-challenge/config"
	server "arena-backend-challenge/internal"
)

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
