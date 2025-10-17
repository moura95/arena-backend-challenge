package repository

import (
	"encoding/csv"
	"fmt"
	"os"
	"sort"
	"strconv"

	"arena-backend-challenge/internal/domain"
)

type MemoryRepository struct {
	locations []domain.Location
}

func NewMemoryRepository(csvPath string) (*MemoryRepository, error) {
	locations, err := loadCSV(csvPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load CSV: %w", err)
	}

	sort.Slice(locations, func(i, j int) bool {
		return locations[i].LowerIPID < locations[j].LowerIPID
	})

	return &MemoryRepository{
		locations: locations,
	}, nil
}

func (r *MemoryRepository) FindByIPID(ipID uint32) (*domain.Location, error) {
	idx := sort.Search(len(r.locations), func(i int) bool {
		return r.locations[i].UpperIPID >= ipID
	})

	if idx < len(r.locations) && r.locations[idx].LowerIPID <= ipID && ipID <= r.locations[idx].UpperIPID {
		return &r.locations[idx], nil
	}

	return nil, domain.ErrLocationNotFound
}

func loadCSV(csvPath string) ([]domain.Location, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, fmt.Errorf("CSV file is empty")
	}

	locations := make([]domain.Location, 0, len(records)-1)

	for i, record := range records {
		if i == 0 {
			continue
		}

		if len(record) < 10 {
			continue
		}

		lowerIPID, err := strconv.ParseUint(record[0], 10, 32)
		if err != nil {
			continue
		}

		upperIPID, err := strconv.ParseUint(record[1], 10, 32)
		if err != nil {
			continue
		}

		location := domain.Location{
			LowerIPID:   uint32(lowerIPID),
			UpperIPID:   uint32(upperIPID),
			Country:     record[2],
			CountryCode: record[3],
			Region:      record[4],
			City:        record[5],
		}

		locations = append(locations, location)
	}

	return locations, nil
}
