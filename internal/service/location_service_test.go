package service

import (
	"errors"
	"testing"

	"arena-backend-challenge/internal/domain"
	"arena-backend-challenge/internal/repository"
)

func TestLocationService_GetLocationByIP(t *testing.T) {
	tests := []struct {
		name          string
		ip            string
		mockFunc      func(ipID uint32) (*domain.Location, error)
		wantLocation  *domain.Location
		wantErr       bool
		expectedError error
	}{
		{
			name: "valid IP - location found",
			ip:   "8.8.8.8",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return &domain.Location{
					LowerIPID:   134744072,
					UpperIPID:   134744072,
					Country:     "United States",
					CountryCode: "US",
					Region:      "California",
					City:        "Mountain View",
				}, nil
			},
			wantLocation: &domain.Location{
				LowerIPID:   134744072,
				UpperIPID:   134744072,
				Country:     "United States",
				CountryCode: "US",
				Region:      "California",
				City:        "Mountain View",
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "valid IP - location not found",
			ip:   "192.168.1.1",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, domain.ErrLocationNotFound
			},
			wantLocation:  nil,
			wantErr:       true,
			expectedError: domain.ErrLocationNotFound,
		},
		{
			name: "invalid IP format",
			ip:   "invalid.ip.address",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, nil
			},
			wantLocation:  nil,
			wantErr:       true,
			expectedError: nil,
		},
		{
			name: "empty IP",
			ip:   "",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, nil
			},
			wantLocation:  nil,
			wantErr:       true,
			expectedError: nil,
		},
		{
			name: "IP with spaces - trimmed correctly",
			ip:   "  8.8.8.8  ",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return &domain.Location{
					Country:     "United States",
					CountryCode: "US",
					Region:      "California",
					City:        "Mountain View",
				}, nil
			},
			wantLocation: &domain.Location{
				Country:     "United States",
				CountryCode: "US",
				Region:      "California",
				City:        "Mountain View",
			},
			wantErr:       false,
			expectedError: nil,
		},
		{
			name: "repository returns unexpected error",
			ip:   "8.8.8.8",
			mockFunc: func(ipID uint32) (*domain.Location, error) {
				return nil, errors.New("database connection error")
			},
			wantLocation:  nil,
			wantErr:       true,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository
			mockRepo := &repository.MockRepository{
				FindByIPIDFunc: tt.mockFunc,
			}

			// Create service with mock
			service := NewLocationService(mockRepo)

			// Call the method
			got, err := service.GetLocationByIP(tt.ip)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("GetLocationByIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check specific error if expected
			if tt.expectedError != nil && !errors.Is(err, tt.expectedError) {
				t.Errorf("GetLocationByIP() error = %v, expectedError %v", err, tt.expectedError)
				return
			}

			// Check location result
			if !tt.wantErr && got != nil && tt.wantLocation != nil {
				if got.Country != tt.wantLocation.Country ||
					got.CountryCode != tt.wantLocation.CountryCode ||
					got.Region != tt.wantLocation.Region ||
					got.City != tt.wantLocation.City {
					t.Errorf("GetLocationByIP() = %+v, want %+v", got, tt.wantLocation)
				}
			}
		})
	}
}

func BenchmarkLocationService_GetLocationByIP(b *testing.B) {
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

	service := NewLocationService(mockRepo)
	testIP := "8.8.8.8"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.GetLocationByIP(testIP)
	}
}
