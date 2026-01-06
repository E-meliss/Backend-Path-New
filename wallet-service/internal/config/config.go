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
	RedisAddr            string
	LogLevel             string
	JWTSecret            string
	AccessTokenTTLMin    int
	RefreshTokenTTLHours int

	DefaultCurrency      string
	DailyDebitLimitCents int64
}

func Load() Config {
	cfg := Config{
		Env:                  getEnv("APP_ENV", "dev"),
		HTTPAddr:             getEnv("HTTP_ADDR", ":8080"),
		DatabaseURL:          getEnv("DATABASE_URL", ""),
		RedisAddr:            getEnv("REDIS_ADDR", ""),
		LogLevel:             getEnv("LOG_LEVEL", "info"),
		JWTSecret:            getEnv("JWT_SECRET", "change-me"),
		AccessTokenTTLMin:    getEnvInt("ACCESS_TOKEN_TTL_MIN", 15),
		RefreshTokenTTLHours: getEnvInt("REFRESH_TOKEN_TTL_HOURS", 24*7),

		DefaultCurrency:      getEnv("DEFAULT_CURRENCY", "USD"),
		DailyDebitLimitCents: getEnvInt64("DAILY_DEBIT_LIMIT_CENTS", 0), // 0 = disabled
	}

	if cfg.DatabaseURL == "" {
		log.Println("WARNING: DATABASE_URL is empty")
	}

	return cfg
}

func getEnv(key, def string) string {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	return v
}

func getEnvInt(key string, def int) int {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return i
}

func getEnvInt64(key string, def int64) int64 {
	v := os.Getenv(key)
	if v == "" {
		return def
	}
	i, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		return def
	}
	return i
}
