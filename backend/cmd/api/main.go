package main

import (
	"errors"
	"log"
	"net/http"

	"pantrypal/backend/internal/app"
	"pantrypal/backend/internal/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal("invalid config: TOKEN_TTL_HOURS must be a positive integer")
	}

	if err := app.Run(cfg); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}
