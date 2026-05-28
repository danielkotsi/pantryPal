package services

import (
	"context"
	"database/sql"
	"errors"

	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

var ErrRecipeNotFound = errors.New("recipe not found")

type RecipeService struct {
	recipes *repositories.RecipeRepository
}

func NewRecipeService(recipes *repositories.RecipeRepository) *RecipeService {
	return &RecipeService{recipes: recipes}
}

func (s *RecipeService) GetRecipe(ctx context.Context, recipeID string) (dto.RecipeResponse, error) {
	recipe, err := s.recipes.GetByID(ctx, recipeID)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.RecipeResponse{}, ErrRecipeNotFound
	}
	return recipe, err
}
