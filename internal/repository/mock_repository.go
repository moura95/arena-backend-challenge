package repository

import "arena-backend-challenge/internal/domain"

type MockRepository struct {
	FindByIPIDFunc func(ipID uint32) (*domain.Location, error)
}

func (m *MockRepository) FindByIPID(ipID uint32) (*domain.Location, error) {
	if m.FindByIPIDFunc != nil {
		return m.FindByIPIDFunc(ipID)
	}
	return nil, domain.ErrLocationNotFound
}
