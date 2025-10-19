package domain

import "errors"

type Location struct {
	LowerIPID   uint32
	UpperIPID   uint32
	Country     string
	CountryCode string
	City        string
}

type Repository interface {
	FindByIPID(ipID uint32) (*Location, error)
}

var (
	ErrLocationNotFound = errors.New("location not found for the given IP")
)
