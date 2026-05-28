package handlers

import (
	"errors"
	"net/http"

	"pantrypal/backend/internal/services"
)

type RecipeHandler struct {
	recipes *services.RecipeService
}

func NewRecipeHandler(recipes *services.RecipeService) *RecipeHandler {
	return &RecipeHandler{recipes: recipes}
}

func (h *RecipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	recipe, err := h.recipes.GetRecipe(r.Context(), r.PathValue("id"))
	if err != nil {
		if errors.Is(err, services.ErrRecipeNotFound) {
			WriteAPIError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to load recipe")
		return
	}

	WriteJSON(w, http.StatusOK, recipe)
}
