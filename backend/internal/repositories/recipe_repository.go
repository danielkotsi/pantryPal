package repositories

import (
	"context"
	"database/sql"

	"pantrypal/backend/internal/transport/http/dto"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) GetByID(ctx context.Context, recipeID string) (dto.RecipeResponse, error) {
	var recipe dto.RecipeResponse
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, name, meal_type, servings, COALESCE(instructions, ''), estimated_cost_cents,
		        total_kcal, total_protein_g, total_carbs_g, total_fat_g
		 FROM recipes
		 WHERE id = ?`,
		recipeID,
	).Scan(
		&recipe.ID,
		&recipe.Name,
		&recipe.MealType,
		&recipe.Servings,
		&recipe.Instructions,
		&recipe.EstimatedCostCents,
		&recipe.Macros.Calories,
		&recipe.Macros.ProteinG,
		&recipe.Macros.CarbsG,
		&recipe.Macros.FatG,
	)
	if err != nil {
		return dto.RecipeResponse{}, err
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT rf.fdc_id, uf.description, rf.quantity, rf.unit
		 FROM recipe_ingredients rf
		 JOIN usda_foods uf ON uf.fdc_id = rf.fdc_id
		 WHERE rf.recipe_id = ?
		 ORDER BY uf.description ASC`,
		recipeID,
	)
	if err != nil {
		return dto.RecipeResponse{}, err
	}
	defer rows.Close()

	recipe.Ingredients = make([]dto.RecipeIngredientResponse, 0)
	for rows.Next() {
		var ingredient dto.RecipeIngredientResponse
		if err := rows.Scan(&ingredient.FDCID, &ingredient.Description, &ingredient.Quantity, &ingredient.Unit); err != nil {
			return dto.RecipeResponse{}, err
		}
		recipe.Ingredients = append(recipe.Ingredients, ingredient)
	}

	return recipe, rows.Err()
}
