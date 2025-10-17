package server

import (
	"fmt"
	"log"
	"net/http"

	"arena-backend-challenge/config"
	"arena-backend-challenge/internal/handler"
	"arena-backend-challenge/internal/repository"
	"arena-backend-challenge/internal/service"
)

type Server struct {
	config          *config.Config
	locationHandler *handler.LocationHandler
}

func NewServer(cfg *config.Config) (*Server, error) {
	repo, err := repository.NewMemoryRepository(cfg.CSVFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	locationService := service.NewLocationService(repo)
	locationHandler := handler.NewLocationHandler(locationService)

	return &Server{
		config:          cfg,
		locationHandler: locationHandler,
	}, nil
}

func (s *Server) Start() error {
	s.registerRoutes()

	log.Printf("Server starting on %s", s.config.HTTPServerAddress)
	return http.ListenAndServe(s.config.HTTPServerAddress, nil)
}

func (s *Server) registerRoutes() {
	http.HandleFunc("/ip/location", s.locationHandler.GetLocation)
}
