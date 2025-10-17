package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPServerAddress string
	CSVFilePath       string
}

func Load() (*Config, error) {
	cfg := &Config{
		HTTPServerAddress: getEnv("HTTP_SERVER_ADDRESS", "0.0.0.0:8080"),
		CSVFilePath:       getEnv("CSV_FILE_PATH", "data/sample.csv"),
	}

	if err := cfg.validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func (c *Config) validate() error {
	if c.HTTPServerAddress == "" {
		return fmt.Errorf("HTTP_SERVER_ADDRESS cannot be empty")
	}
	if c.CSVFilePath == "" {
		return fmt.Errorf("CSV_FILE_PATH cannot be empty")
	}
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
