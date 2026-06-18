package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all runtime configuration, loaded from environment variables.
type Config struct {
	DatabaseURL string
	JWTSecret   string
	JWTTTLHours int
	Port        string
	CORSOrigins []string
}

// Load reads configuration from the environment, applying sensible defaults for
// local development and returning an error only when a required secret is missing.
func Load() (Config, error) {
	cfg := Config{
		DatabaseURL: getenv("DATABASE_URL", "postgres://fitness:fitness@localhost:5432/fitness?sslmode=disable"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTTTLHours: getenvInt("JWT_TTL_HOURS", 720),
		Port:        getenv("PORT", "8080"),
		CORSOrigins: splitCSV(getenv("CORS_ORIGINS", "http://localhost:3000")),
	}

	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET is required")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getenvInt(key string, fallback int) int {
	if v := os.Getenv(key); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func splitCSV(s string) []string {
	parts := strings.Split(s, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			out = append(out, t)
		}
	}
	return out
}
