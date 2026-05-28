package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port                 string
	DBPath               string
	TokenSecret          string
	TokenTTL             time.Duration
	GeminiAPIKey         string
	GeminiModel          string
	GeminiBaseURL        string
	GeminiTimeout        time.Duration
	GeminiRetryMax       int
	GeminiRetryBackoff   time.Duration
	GeminiResponseFormat string
}

func Load() (Config, error) {
	tokenHours, err := strconv.Atoi(envOr("TOKEN_TTL_HOURS", "24"))
	if err != nil {
		return Config{}, err
	}
	if tokenHours < 1 {
		return Config{}, errors.New("TOKEN_TTL_HOURS must be positive")
	}

	geminiTimeoutSeconds, err := strconv.Atoi(envOr("GEMINI_TIMEOUT_SECONDS", "20"))
	if err != nil {
		return Config{}, err
	}
	if geminiTimeoutSeconds < 1 {
		return Config{}, errors.New("GEMINI_TIMEOUT_SECONDS must be positive")
	}

	geminiRetryMax, err := strconv.Atoi(envOr("GEMINI_RETRY_MAX", "2"))
	if err != nil {
		return Config{}, err
	}
	if geminiRetryMax < 0 {
		return Config{}, errors.New("GEMINI_RETRY_MAX must be zero or positive")
	}

	geminiRetryBackoffMS, err := strconv.Atoi(envOr("GEMINI_RETRY_BACKOFF_MS", "500"))
	if err != nil {
		return Config{}, err
	}
	if geminiRetryBackoffMS < 1 {
		return Config{}, errors.New("GEMINI_RETRY_BACKOFF_MS must be positive")
	}

	return Config{
		Port:                 envOr("PORT", "8080"),
		DBPath:               envOr("DB_PATH", "../database/sqlite/pantrypal.db"),
		TokenSecret:          envOr("TOKEN_SECRET", "pantrypal-local-dev-secret"),
		TokenTTL:             time.Duration(tokenHours) * time.Hour,
		GeminiAPIKey:         strings.TrimSpace(os.Getenv("GEMINI_API_KEY")),
		GeminiModel:          envOr("GEMINI_MODEL", "gemini-1.5-flash"),
		GeminiBaseURL:        envOr("GEMINI_BASE_URL", "https://generativelanguage.googleapis.com"),
		GeminiTimeout:        time.Duration(geminiTimeoutSeconds) * time.Second,
		GeminiRetryMax:       geminiRetryMax,
		GeminiRetryBackoff:   time.Duration(geminiRetryBackoffMS) * time.Millisecond,
		GeminiResponseFormat: envOr("GEMINI_RESPONSE_FORMAT", "application/json"),
	}, nil
}

func envOr(name, fallback string) string {
	v := strings.TrimSpace(os.Getenv(name))
	if v == "" {
		return fallback
	}
	return v
}
