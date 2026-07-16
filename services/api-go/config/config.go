package config

import (
	"os"
)

type Config struct {
	DatabaseURL         string
	RedisURL            string
	CryptoServiceURL    string
	JWTSecret           string
	Port                string
	RateLimitRPM        int
	AccessTokenExpiry   int // minutes
	RefreshTokenExpiry  int // hours
}

func Load() *Config {
	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", "postgres://onentry:onentry_password@localhost:5432/onentry?sslmode=disable"),
		RedisURL:           getEnv("REDIS_URL", "redis://localhost:6379"),
		CryptoServiceURL:   getEnv("CRYPTO_SERVICE_URL", "http://localhost:8083"),
		JWTSecret:          getEnv("JWT_SECRET", "change-me-in-production"),
		Port:               getEnv("PORT", "8082"),
		RateLimitRPM:       60,
		AccessTokenExpiry:  15,
		RefreshTokenExpiry: 720,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
