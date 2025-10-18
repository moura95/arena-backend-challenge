package repository

import (
	"os"
	"path/filepath"
	"testing"

	"arena-backend-challenge/internal/domain"
)

func TestNewMemoryRepository(t *testing.T) {
	tests := []struct {
		name      string
		csvData   string
		wantErr   bool
		wantCount int
	}{
		{
			name: "valid CSV with multiple rows",
			csvData: `"ip_from","ip_to","country_code","country_name","region_name","city_name","latitude","longitude","zip_code","time_zone"
"16777216","16777471","US","United States","California","Los Angeles","34.05223","-118.24368","90001","-07:00"
"16777472","16778239","CN","China","Fujian","Fuzhou","26.06139","119.30611","-","08:00"
"16778240","16779263","AU","Australia","Queensland","Brisbane","-27.46794","153.02809","4000","10:00"`,
			wantErr:   false,
			wantCount: 3,
		},
		{
			name:      "empty CSV",
			csvData:   "",
			wantErr:   true,
			wantCount: 0,
		},
		{
			name:      "CSV with only header",
			csvData:   `"ip_from","ip_to","country_code","country_name","region_name","city_name","latitude","longitude","zip_code","time_zone"`,
			wantErr:   false,
			wantCount: 0,
		},
		{
			name: "CSV with invalid rows (skipped)",
			csvData: `"ip_from","ip_to","country_code","country_name","region_name","city_name","latitude","longitude","zip_code","time_zone"
"16777216","16777471","US","United States","California","Los Angeles","34.05223","-118.24368","90001","-07:00"
"invalid","16778239","CN","China","Fujian","Fuzhou","26.06139","119.30611","-","08:00"
"16778240","16779263","AU","Australia","Queensland","Brisbane","-27.46794","153.02809","4000","10:00"`,
			wantErr:   false,
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary CSV file
			tmpFile, err := createTempCSV(tt.csvData)
			if err != nil {
				t.Fatalf("Failed to create temp CSV: %v", err)
			}
			defer os.Remove(tmpFile)

			// Create repository
			repo, err := NewMemoryRepository(tmpFile)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("NewMemoryRepository() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check count if no error
			if !tt.wantErr && repo != nil {
				if len(repo.locations) != tt.wantCount {
					t.Errorf("NewMemoryRepository() loaded %d locations, want %d", len(repo.locations), tt.wantCount)
				}
			}
		})
	}
}

func TestMemoryRepository_FindByIPID(t *testing.T) {
	// Create test CSV
	csvData := `"ip_from","ip_to","country_code","country_name","region_name","city_name","latitude","longitude","zip_code","time_zone"
"16777216","16777471","US","United States","California","Los Angeles","34.05223","-118.24368","90001","-07:00"
"16777472","16778239","CN","China","Fujian","Fuzhou","26.06139","119.30611","-","08:00"
"16778240","16779263","AU","Australia","Queensland","Brisbane","-27.46794","153.02809","4000","10:00"
"134744072","134744072","US","United States","California","Mountain View","37.405992","-122.078515","94035","-07:00"`

	tmpFile, err := createTempCSV(csvData)
	if err != nil {
		t.Fatalf("Failed to create temp CSV: %v", err)
	}
	defer os.Remove(tmpFile)

	repo, err := NewMemoryRepository(tmpFile)
	if err != nil {
		t.Fatalf("Failed to create repository: %v", err)
	}

	tests := []struct {
		name        string
		ipID        uint32
		wantCountry string
		wantCity    string
		wantErr     bool
		expectedErr error
	}{
		{
			name:        "IP in first range - lower bound",
			ipID:        16777216,
			wantCountry: "United States",
			wantCity:    "Los Angeles",
			wantErr:     false,
		},
		{
			name:        "IP in first range - upper bound",
			ipID:        16777471,
			wantCountry: "United States",
			wantCity:    "Los Angeles",
			wantErr:     false,
		},
		{
			name:        "IP in first range - middle",
			ipID:        16777350,
			wantCountry: "United States",
			wantCity:    "Los Angeles",
			wantErr:     false,
		},
		{
			name:        "IP in second range",
			ipID:        16777500,
			wantCountry: "China",
			wantCity:    "Fuzhou",
			wantErr:     false,
		},
		{
			name:        "IP in third range",
			ipID:        16778500,
			wantCountry: "Australia",
			wantCity:    "Brisbane",
			wantErr:     false,
		},
		{
			name:        "Exact IP match (8.8.8.8)",
			ipID:        134744072,
			wantCountry: "United States",
			wantCity:    "Mountain View",
			wantErr:     false,
		},
		{
			name:        "IP not in any range - below all",
			ipID:        1000,
			wantErr:     true,
			expectedErr: domain.ErrLocationNotFound,
		},
		{
			name:        "IP not in any range - above all",
			ipID:        999999999,
			wantErr:     true,
			expectedErr: domain.ErrLocationNotFound,
		},
		{
			name:        "IP in gap between ranges",
			ipID:        16779300,
			wantErr:     true,
			expectedErr: domain.ErrLocationNotFound,
		},
		{
			name:        "IP at zero",
			ipID:        0,
			wantErr:     true,
			expectedErr: domain.ErrLocationNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindByIPID(tt.ipID)

			// Check error
			if (err != nil) != tt.wantErr {
				t.Errorf("FindByIPID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// Check specific error
			if tt.expectedErr != nil && err != tt.expectedErr {
				t.Errorf("FindByIPID() error = %v, expectedErr %v", err, tt.expectedErr)
				return
			}

			// Check result if no error
			if !tt.wantErr && got != nil {
				if got.Country != tt.wantCountry {
					t.Errorf("FindByIPID() Country = %v, want %v", got.Country, tt.wantCountry)
				}
				if got.City != tt.wantCity {
					t.Errorf("FindByIPID() City = %v, want %v", got.City, tt.wantCity)
				}
			}
		})
	}
}

func BenchmarkMemoryRepository_FindByIPID(b *testing.B) {
	// Create test CSV with many rows for realistic benchmark
	csvData := `"ip_from","ip_to","country_code","country_name","region_name","city_name","latitude","longitude","zip_code","time_zone"
"16777216","16777471","US","United States","California","Los Angeles","34.05223","-118.24368","90001","-07:00"
"16777472","16778239","CN","China","Fujian","Fuzhou","26.06139","119.30611","-","08:00"
"16778240","16779263","AU","Australia","Queensland","Brisbane","-27.46794","153.02809","4000","10:00"
"134744072","134744072","US","United States","California","Mountain View","37.405992","-122.078515","94035","-07:00"
"167772160","167772415","US","United States","New York","New York","40.7127","-74.0059","10001","-05:00"`

	tmpFile, err := createTempCSV(csvData)
	if err != nil {
		b.Fatalf("Failed to create temp CSV: %v", err)
	}
	defer os.Remove(tmpFile)

	repo, err := NewMemoryRepository(tmpFile)
	if err != nil {
		b.Fatalf("Failed to create repository: %v", err)
	}

	testIPID := uint32(134744072) // 8.8.8.8

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = repo.FindByIPID(testIPID)
	}
}

func createTempCSV(content string) (string, error) {
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "test_ip_data.csv")

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		return "", err
	}

	return tmpFile, nil
}
