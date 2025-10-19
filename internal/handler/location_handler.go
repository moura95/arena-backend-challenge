package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	v1 "arena-backend-challenge/api/v1"
	"arena-backend-challenge/internal/domain"
	"arena-backend-challenge/internal/service"
	"arena-backend-challenge/pkg/logger"
)

type LocationHandler struct {
	service *service.LocationService
}

func NewLocationHandler(service *service.LocationService) *LocationHandler {
	return &LocationHandler{
		service: service,
	}
}

func (h *LocationHandler) GetLocation(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	ip := r.URL.Query().Get("ip")
	if ip == "" {
		h.sendError(w, "IP address is required", http.StatusBadRequest)
		logger.Warningf("Bad request - missing IP parameter - Duration: %v", time.Since(start))
		return
	}

	location, err := h.service.GetLocationByIP(ip)
	if err != nil {
		if errors.Is(err, domain.ErrLocationNotFound) {
			h.sendError(w, "Location not found for the given IP", http.StatusNotFound)
			logger.Infof("IP lookup - IP: %s - Status: 404 Not Found - Duration: %v", ip, time.Since(start))
			return
		}
		h.sendError(w, err.Error(), http.StatusBadRequest)
		logger.Warningf("Invalid IP - IP: %s - Error: %v - Duration: %v", ip, err, time.Since(start))
		return
	}

	response := v1.LocationResponse{
		Country:     location.Country,
		CountryCode: location.CountryCode,
		Region:      location.Region,
		City:        location.City,
	}

	h.sendJSON(w, response, http.StatusOK)
	logger.Infof("IP lookup - IP: %s - Country: %s - City: %s - Status: 200 OK - Duration: %v",
		ip, location.Country, location.City, time.Since(start))
}

func (h *LocationHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logger.Errorf("Error encoding JSON response: %v", err)
	}
}

func (h *LocationHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	h.sendJSON(w, v1.ErrorResponse{Error: message}, statusCode)
}
