package config

import (
	"log"
	"os"
	"strconv"
)

type Config struct {
	Env                  string
	HTTPAddr             string
	DatabaseURL          string
	LogLevel             string
	JWTSecret            string
	AccessTokenTTLMin    int
	RefreshTokenTTLHours int
}

func Load() Config {
	cfg := Config{
		Env:                  getEnv("APP_ENV", "dev"),
		HTTPAddr:             getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:          mustEnv("DATABASE_URL"),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		JWTSecret:            mustEnv("JWT_SECRET"),
		AccessTokenTTLMin:    mustInt("ACCESS_TOKEN_TTL_MIN"),
		RefreshTokenTTLHours: mustInt("REFRESH_TOKEN_TTL_HOURS"),
	}
	return cfg
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func mustEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("missing required env: %s", key)
	}
	return v
}

func mustInt(key string) int {
	v := mustEnv(key)
	n, err := strconv.Atoi(v)
	if err != nil {
		log.Fatalf("invalid int env %s=%s", key, v)
	}
	return n
}
