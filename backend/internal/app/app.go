package app

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"pantrypal/backend/internal/config"
	"pantrypal/backend/internal/platform/auth"
	"pantrypal/backend/internal/platform/db"
	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/handlers"
	"pantrypal/backend/internal/transport/http/router"
)

func Run(cfg config.Config) error {
	conn, err := db.Open(cfg.DBPath)
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	defer conn.Close()

	userRepo := repositories.NewUserRepository(conn)
	foodRepo := repositories.NewFoodRepository(conn)
	pantryRepo := repositories.NewPantryRepository(conn)
	recipeRepo := repositories.NewRecipeRepository(conn)
	planRepo := repositories.NewPlanRepository(conn)
	tokenManager := auth.NewTokenManager(cfg.TokenSecret, cfg.TokenTTL)

	authService := services.NewAuthService(userRepo, tokenManager)
	profileService := services.NewProfileService(userRepo)
	pantryService := services.NewPantryService(foodRepo, pantryRepo)
	recipeService := services.NewRecipeService(recipeRepo)
	planService := services.NewPlanService(planRepo)

	healthHandler := handlers.NewHealthHandler()
	authHandler := handlers.NewAuthHandler(authService)
	profileHandler := handlers.NewProfileHandler(profileService)
	pantryHandler := handlers.NewPantryHandler(pantryService)
	recipeHandler := handlers.NewRecipeHandler(recipeService)
	planHandler := handlers.NewPlanHandler(planService)

	rootHandler := router.New(router.Handlers{
		Health:  healthHandler,
		Auth:    authHandler,
		Profile: profileHandler,
		Pantry:  pantryHandler,
		Recipe:  recipeHandler,
		Plan:    planHandler,
	}, tokenManager, userRepo)

	server := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           rootHandler,
		ReadHeaderTimeout: 10 * time.Second,
	}

	log.Printf("api listening on :%s using db %s", cfg.Port, cfg.DBPath)
	return server.ListenAndServe()
}
