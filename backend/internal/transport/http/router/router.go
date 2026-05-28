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
	Pantry  *handlers.PantryHandler
	Recipe  *handlers.RecipeHandler
	Plan    *handlers.PlanHandler
}

func New(h Handlers, tokens *auth.TokenManager, users *repositories.UserRepository) http.Handler {
	mux := http.NewServeMux()

	authRequired := middleware.AuthRequired(tokens, users)

	mux.HandleFunc("GET /health", h.Health.GetHealth)
	mux.HandleFunc("POST /auth/register", h.Auth.Register)
	mux.HandleFunc("POST /auth/login", h.Auth.Login)
	mux.Handle("GET /ingredients/search", authRequired(http.HandlerFunc(h.Pantry.SearchFoods)))
	mux.Handle("GET /me", authRequired(http.HandlerFunc(h.Auth.Me)))
	mux.Handle("GET /profile", authRequired(http.HandlerFunc(h.Profile.GetProfile)))
	mux.Handle("PATCH /profile/metrics", authRequired(http.HandlerFunc(h.Profile.PatchMetrics)))
	mux.Handle("PATCH /profile/preferences", authRequired(http.HandlerFunc(h.Profile.PatchPreferences)))
	mux.Handle("PATCH /profile/budget", authRequired(http.HandlerFunc(h.Profile.PatchBudget)))
	mux.Handle("GET /pantry/items", authRequired(http.HandlerFunc(h.Pantry.ListPantryItems)))
	mux.Handle("POST /pantry/items", authRequired(http.HandlerFunc(h.Pantry.CreatePantryItem)))
	mux.Handle("PATCH /pantry/items/{id}", authRequired(http.HandlerFunc(h.Pantry.PatchPantryItem)))
	mux.Handle("DELETE /pantry/items/{id}", authRequired(http.HandlerFunc(h.Pantry.DeletePantryItem)))
	mux.Handle("GET /recipes/{id}", authRequired(http.HandlerFunc(h.Recipe.GetRecipe)))
	mux.Handle("POST /plans/proposal", authRequired(http.HandlerFunc(h.Plan.CreateProposal)))
	mux.Handle("POST /plans/{id}/accept", authRequired(http.HandlerFunc(h.Plan.AcceptProposal)))
	mux.Handle("POST /plans/{id}/decline", authRequired(http.HandlerFunc(h.Plan.DeclineProposal)))
	mux.Handle("GET /plans/week", authRequired(http.HandlerFunc(h.Plan.GetWeekPlan)))

	return middleware.CORS(middleware.Logging(mux))
}
