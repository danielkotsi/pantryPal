package config

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	Port        string
	DBPath      string
	TokenSecret string
	TokenTTL    time.Duration
}

func Load() (Config, error) {
	tokenHours, err := strconv.Atoi(envOr("TOKEN_TTL_HOURS", "24"))
	if err != nil {
		return Config{}, err
	}
	if tokenHours < 1 {
		return Config{}, errors.New("TOKEN_TTL_HOURS must be positive")
	}

	return Config{
		Port:        envOr("PORT", "8080"),
		DBPath:      envOr("DB_PATH", "../database/sqlite/pantrypal.db"),
		TokenSecret: envOr("TOKEN_SECRET", "pantrypal-local-dev-secret"),
		TokenTTL:    time.Duration(tokenHours) * time.Hour,
	}, nil
}

func envOr(name, fallback string) string {
	v := strings.TrimSpace(os.Getenv(name))
	if v == "" {
		return fallback
	}
	return v
}
