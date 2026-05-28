package handlers

import (
	"errors"
	"net/http"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/middleware"
)

type ConsumeHandler struct {
	consume *services.ConsumeService
}

func NewConsumeHandler(consume *services.ConsumeService) *ConsumeHandler {
	return &ConsumeHandler{consume: consume}
}

func (h *ConsumeHandler) ConsumeMeal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	mealID := r.PathValue("id")
	if mealID == "" {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "meal id is required")
		return
	}

	res, err := h.consume.ConsumeMeal(r.Context(), userID, mealID)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrPlanMealNotFound):
			WriteAPIError(w, http.StatusNotFound, "not_found", err.Error())
		case errors.Is(err, services.ErrMealAlreadyConsumed):
			WriteAPIError(w, http.StatusConflict, "conflict", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to consume meal")
		}
		return
	}

	WriteJSON(w, http.StatusOK, res)
}
