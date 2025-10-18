package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	v1 "arena-backend-challenge/api/v1"
	"arena-backend-challenge/internal/domain"
	"arena-backend-challenge/internal/repository"
	"arena-backend-challenge/internal/service"
)

func TestLocationHandler_GetLocation(t *testing.T) {
	tests := []struct {
		name           string
		queryParam     string
		mockFunc       func(ipID uint32) (*domain.Location, error)
		wantStatus     int
		wantCountry    string
		wantCity       string
		wantError      string
		checkErrorOnly bool
	}{
		{
			name:       "valid IP - location found",
			queryParam: "ip=8.8.8.8",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return &domain.Location{
					Country:     "United States",
					CountryCode: "US",
					Region:      "California",
					City:        "Mountain View",
				}, nil
			},
			wantStatus:  http.StatusOK,
			wantCountry: "United States",
			wantCity:    "Mountain View",
		},
		{
			name:       "valid IP - location not found",
			queryParam: "ip=192.168.1.1",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, domain.ErrLocationNotFound
			},
			wantStatus:     http.StatusNotFound,
			wantError:      "Location not found for the given IP",
			checkErrorOnly: true,
		},
		{
			name:           "missing IP parameter",
			queryParam:     "",
			mockFunc:       nil,
			wantStatus:     http.StatusBadRequest,
			wantError:      "IP address is required",
			checkErrorOnly: true,
		},
		{
			name:       "invalid IP format",
			queryParam: "ip=invalid.ip.address",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, nil
			},
			wantStatus:     http.StatusBadRequest,
			checkErrorOnly: true,
		},
		{
			name:       "IP with spaces",
			queryParam: "ip=8.8.8.8",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return &domain.Location{
					Country:     "United States",
					CountryCode: "US",
					Region:      "California",
					City:        "Mountain View",
				}, nil
			},
			wantStatus:  http.StatusOK,
			wantCountry: "United States",
			wantCity:    "Mountain View",
		},
		{
			name:       "IP out of range",
			queryParam: "ip=256.256.256.256",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, nil
			},
			wantStatus:     http.StatusBadRequest,
			checkErrorOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			var mockRepo *repository.MockRepository
			if tt.mockFunc != nil {
				mockRepo = &repository.MockRepository{
					FindByIPIDFunc: tt.mockFunc,
				}
			} else {
				mockRepo = &repository.MockRepository{}
			}

			// Create service and handler
			locationService := service.NewLocationService(mockRepo)
			handler := NewLocationHandler(locationService)

			// Create request
			req := httptest.NewRequest(http.MethodGet, "/ip/location?"+tt.queryParam, nil)
			w := httptest.NewRecorder()

			// Call handler
			handler.GetLocation(w, req)

			// Check status code
			if w.Code != tt.wantStatus {
				t.Errorf("GetLocation() status = %v, want %v", w.Code, tt.wantStatus)
			}

			// Check Content-Type
			contentType := w.Header().Get("Content-Type")
			if contentType != "application/json" {
				t.Errorf("GetLocation() Content-Type = %v, want application/json", contentType)
			}

			// Parse response
			if tt.checkErrorOnly {
				// Check error response
				var errorResp v1.ErrorResponse
				if err := json.NewDecoder(w.Body).Decode(&errorResp); err != nil {
					t.Fatalf("Failed to decode error response: %v", err)
				}
				if tt.wantError != "" && errorResp.Error != tt.wantError {
					t.Errorf("GetLocation() error = %v, want %v", errorResp.Error, tt.wantError)
				}
			} else {
				// Check success response
				var locationResp v1.LocationResponse
				if err := json.NewDecoder(w.Body).Decode(&locationResp); err != nil {
					t.Fatalf("Failed to decode location response: %v", err)
				}
				if locationResp.Country != tt.wantCountry {
					t.Errorf("GetLocation() Country = %v, want %v", locationResp.Country, tt.wantCountry)
				}
				if locationResp.City != tt.wantCity {
					t.Errorf("GetLocation() City = %v, want %v", locationResp.City, tt.wantCity)
				}
			}
		})
	}
}

func TestLocationHandler_GetLocation_Methods(t *testing.T) {
	mockRepo := &repository.MockRepository{
		FindByIPIDFunc: func(ipID uint32) (*domain.Location, error) {
			return &domain.Location{
				Country:     "United States",
				CountryCode: "US",
			}, nil
		},
	}

	locationService := service.NewLocationService(mockRepo)
	handler := NewLocationHandler(locationService)

	tests := []struct {
		name       string
		method     string
		wantStatus int
	}{
		{
			name:       "GET method - allowed",
			method:     http.MethodGet,
			wantStatus: http.StatusOK,
		},
		{
			name:       "POST method - should work (net/http doesn't restrict)",
			method:     http.MethodPost,
			wantStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/ip/location?ip=8.8.8.8", nil)
			w := httptest.NewRecorder()

			handler.GetLocation(w, req)

			if w.Code != tt.wantStatus {
				t.Errorf("GetLocation() with %s method status = %v, want %v", tt.method, w.Code, tt.wantStatus)
			}
		})
	}
}

func BenchmarkLocationHandler_GetLocation(b *testing.B) {
	mockRepo := &repository.MockRepository{
		FindByIPIDFunc: func(ipID uint32) (*domain.Location, error) {
			return &domain.Location{
				Country:     "United States",
				CountryCode: "US",
				Region:      "California",
				City:        "Mountain View",
			}, nil
		},
	}

	locationService := service.NewLocationService(mockRepo)
	handler := NewLocationHandler(locationService)

	req := httptest.NewRequest(http.MethodGet, "/ip/location?ip=8.8.8.8", nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		w := httptest.NewRecorder()
		handler.GetLocation(w, req)
	}
}
