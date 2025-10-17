package service

import (
	"fmt"

	"arena-backend-challenge/internal/domain"
	"arena-backend-challenge/pkg/iputil"
)

type LocationService struct {
	repo domain.Repository
}

func NewLocationService(repo domain.Repository) *LocationService {
	return &LocationService{
		repo: repo,
	}
}

func (s *LocationService) GetLocationByIP(ip string) (*domain.Location, error) {
	ipID, err := iputil.IPToID(ip)
	if err != nil {
		return nil, fmt.Errorf("invalid IP address: %w", err)
	}

	location, err := s.repo.FindByIPID(ipID)
	if err != nil {
		return nil, err
	}

	return location, nil
}
