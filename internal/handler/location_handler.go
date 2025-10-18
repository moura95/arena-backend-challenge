package handler

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	v1 "arena-backend-challenge/api/v1"
	"arena-backend-challenge/internal/domain"
	"arena-backend-challenge/internal/service"
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
	ip := r.URL.Query().Get("ip")
	if ip == "" {
		h.sendError(w, "IP address is required", http.StatusBadRequest)
		return
	}

	location, err := h.service.GetLocationByIP(ip)
	if err != nil {
		if errors.Is(err, domain.ErrLocationNotFound) {
			h.sendError(w, "Location not found for the given IP", http.StatusNotFound)
			return
		}
		h.sendError(w, err.Error(), http.StatusBadRequest)
		return
	}

	response := v1.LocationResponse{
		Country:     location.Country,
		CountryCode: location.CountryCode,
		Region:      location.Region,
		City:        location.City,
	}

	h.sendJSON(w, response, http.StatusOK)
}

func (h *LocationHandler) sendJSON(w http.ResponseWriter, data interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
	}
}

func (h *LocationHandler) sendError(w http.ResponseWriter, message string, statusCode int) {
	h.sendJSON(w, v1.ErrorResponse{Error: message}, statusCode)
}
