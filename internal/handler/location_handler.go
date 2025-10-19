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

// GetLocation godoc
// @Summary Get IP location
// @Description Get geographic location information for a given IP address
// @Tags Location
// @Accept json
// @Produce json
// @Param ip query string true "IPv4 address (e.g., 8.8.8.8)"
// @Success 200 {object} v1.LocationResponse "Location found"
// @Failure 400 {object} v1.ErrorResponse "Invalid IP address format"
// @Failure 404 {object} v1.ErrorResponse "Location not found for the given IP"
// @Router /ip/location [get]
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
		duration := time.Since(start)

		if errors.Is(err, domain.ErrLocationNotFound) {
			h.sendError(w, "Location not found for the given IP", http.StatusNotFound)
			logger.Infof("IP lookup not found - IP: %s - Status: 404 - Duration: %v - Error chain: %v",
				ip, duration, err)
			return
		}

		h.sendError(w, err.Error(), http.StatusBadRequest)

		logger.Warningf("IP lookup failed - IP: %s - Status: 400 - Duration: %v - Error chain: %v",
			ip, duration, err)
		return
	}

	response := v1.LocationResponse{
		Country:     location.Country,
		CountryCode: location.CountryCode,
		City:        location.City,
	}

	h.sendJSON(w, response, http.StatusOK)

	duration := time.Since(start)
	logger.Infof("IP lookup success - IP: %s - Country: %s - City: %s - Status: 200 - Duration: %v",
		ip, location.Country, location.City, duration)
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
