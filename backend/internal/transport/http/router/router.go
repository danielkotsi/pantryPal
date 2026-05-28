package router

import (
	"net/http"

	"pantrypal/backend/internal/platform/auth"
	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/handlers"
	"pantrypal/backend/internal/transport/http/middleware"
)

type Handlers struct {
	Health  *handlers.HealthHandler
	Auth    *handlers.AuthHandler
	Profile *handlers.ProfileHandler
}

func New(h Handlers, tokens *auth.TokenManager, users *repositories.UserRepository) http.Handler {
	mux := http.NewServeMux()

	authRequired := middleware.AuthRequired(tokens, users)

	mux.HandleFunc("GET /health", h.Health.GetHealth)
	mux.HandleFunc("POST /auth/register", h.Auth.Register)
	mux.HandleFunc("POST /auth/login", h.Auth.Login)
	mux.Handle("GET /me", authRequired(http.HandlerFunc(h.Auth.Me)))
	mux.Handle("GET /profile", authRequired(http.HandlerFunc(h.Profile.GetProfile)))
	mux.Handle("PATCH /profile/metrics", authRequired(http.HandlerFunc(h.Profile.PatchMetrics)))
	mux.Handle("PATCH /profile/preferences", authRequired(http.HandlerFunc(h.Profile.PatchPreferences)))
	mux.Handle("PATCH /profile/budget", authRequired(http.HandlerFunc(h.Profile.PatchBudget)))

	return middleware.CORS(middleware.Logging(mux))
}
