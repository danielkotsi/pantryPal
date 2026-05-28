package handlers

import (
	"errors"
	"net/http"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/dto"
	"pantrypal/backend/internal/transport/http/middleware"
)

type PantryHandler struct {
	pantry *services.PantryService
}

func NewPantryHandler(pantry *services.PantryService) *PantryHandler {
	return &PantryHandler{pantry: pantry}
}

func (h *PantryHandler) SearchFoods(w http.ResponseWriter, r *http.Request) {
	items, err := h.pantry.SearchFoods(r.Context(), r.URL.Query().Get("q"))
	if err != nil {
		if errors.Is(err, services.ErrInvalidSearchQuery) {
			WriteAPIError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to search foods")
		return
	}

	WriteJSON(w, http.StatusOK, dto.FoodSearchResponse{Items: items})
}

func (h *PantryHandler) ListPantryItems(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	items, err := h.pantry.ListPantryItems(r.Context(), userID)
	if err != nil {
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to load pantry items")
		return
	}

	WriteJSON(w, http.StatusOK, dto.PantryItemsResponse{Items: items})
}

func (h *PantryHandler) CreatePantryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.PantryItemRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	item, err := h.pantry.AddPantryItem(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidQuantity), errors.Is(err, services.ErrInvalidUnit), errors.Is(err, services.ErrInvalidPantryFood):
			WriteAPIError(w, http.StatusBadRequest, "validation_error", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to add pantry item")
		}
		return
	}

	WriteJSON(w, http.StatusCreated, item)
}

func (h *PantryHandler) PatchPantryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.PantryItemPatchRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	item, err := h.pantry.PatchPantryItem(r.Context(), userID, r.PathValue("id"), req)
	if err != nil {
		if errors.Is(err, services.ErrPantryItemNotFound) {
			WriteAPIError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to update pantry item")
		return
	}

	WriteJSON(w, http.StatusOK, item)
}

func (h *PantryHandler) DeletePantryItem(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	err := h.pantry.DeletePantryItem(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		if errors.Is(err, services.ErrPantryItemNotFound) {
			WriteAPIError(w, http.StatusNotFound, "not_found", err.Error())
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to delete pantry item")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
