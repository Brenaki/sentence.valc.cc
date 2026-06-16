package config

import (
	"fmt"
	"os"
)

// Config holds runtime configuration loaded from the environment.
type Config struct {
	HTTPAddr     string
	NinjaAPIKey  string
	MySQLDSN     string
	AllowOrigins string
}

// Load reads configuration from environment variables, applying defaults.
func Load() (*Config, error) {
	cfg := &Config{
		HTTPAddr:     getenv("HTTP_ADDR", ":8080"),
		NinjaAPIKey:  os.Getenv("API_KEY_NINJA"),
		MySQLDSN:     getenv("MYSQL_DSN", "root:root@tcp(127.0.0.1:3306)/sentence?parseTime=true&charset=utf8mb4&loc=Local"),
		AllowOrigins: getenv("ALLOW_ORIGINS", "*"),
	}
	if cfg.NinjaAPIKey == "" {
		return nil, fmt.Errorf("API_KEY_NINJA is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
