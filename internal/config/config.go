package config

import (
	"os"
	"strconv"
)

type Config struct {
	AppEnv    string
	AppPort   string
	SecretKey string

	DBHost         string
	DBPort         string
	DBName         string
	DBUser         string
	DBPassword     string
	DBMaxOpenConns int
	DBMaxIdleConns int

	SessionCookieName  string
	SessionDurationDays int

	EmailFrom     string
	EmailFromName string
	ResendAPIKey  string

	AppBaseURL string
}

func Load() *Config {
	return &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		AppPort:   getEnv("APP_PORT", "8080"),
		SecretKey: getEnv("APP_SECRET_KEY", "dev-secret-change-in-production"),

		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "3306"),
		DBName:         getEnv("DB_NAME", "time_tracker"),
		DBUser:         getEnv("DB_USER", "root"),
		DBPassword:     getEnv("DB_PASSWORD", ""),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 10),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 5),

		SessionCookieName:   getEnv("SESSION_COOKIE_NAME", "tt_session"),
		SessionDurationDays: getEnvInt("SESSION_DURATION_DAYS", 7),

		EmailFrom:     getEnv("EMAIL_FROM", "noreply@localhost"),
		EmailFromName: getEnv("EMAIL_FROM_NAME", "Time Tracker"),
		ResendAPIKey:  getEnv("RESEND_API_KEY", ""),

		AppBaseURL: getEnv("APP_BASE_URL", "http://localhost:8080"),
	}
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if val := os.Getenv(key); val != "" {
		if i, err := strconv.Atoi(val); err == nil {
			return i
		}
	}
	return fallback
}
