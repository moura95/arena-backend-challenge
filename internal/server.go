package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	v1 "arena-backend-challenge/api/v1"
	"arena-backend-challenge/config"
	"arena-backend-challenge/internal/handler"
	"arena-backend-challenge/internal/repository"
	"arena-backend-challenge/internal/service"
	"arena-backend-challenge/pkg/logger"
)

const Version = "1.0.0"

type Server struct {
	config          *config.Config
	locationHandler *handler.LocationHandler
	startTime       time.Time
}

func NewServer(cfg *config.Config) (*Server, error) {
	logger.Info("Initializing server...")

	repo, err := repository.NewMemoryRepository(cfg.CSVFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize repository: %w", err)
	}

	locationService := service.NewLocationService(repo)
	locationHandler := handler.NewLocationHandler(locationService)

	return &Server{
		config:          cfg,
		locationHandler: locationHandler,
		startTime:       time.Now(),
	}, nil
}

func (s *Server) Start() error {
	s.registerRoutes()

	logger.Infof("Server starting on %s (version %s)", s.config.HTTPServerAddress, Version)
	return http.ListenAndServe(s.config.HTTPServerAddress, nil)
}

func (s *Server) registerRoutes() {
	http.HandleFunc("/ip/location", s.locationHandler.GetLocation)
	http.HandleFunc("/health", s.handleHealth)

	logger.Info("Routes registered:")
	logger.Info("  GET /ip/location?ip=<address>")
	logger.Info("  GET /health")
}

func (s *Server) handleHealth(w http.ResponseWriter, r *http.Request) {
	response := v1.HealthResponse{
		Status:    "healthy",
		Timestamp: time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.Errorf("Error encoding health response: %v", err)
	}
}
