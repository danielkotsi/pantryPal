package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/dto"
	"pantrypal/backend/internal/transport/http/middleware"
)

type ProfileHandler struct {
	profile *services.ProfileService
}

func NewProfileHandler(profile *services.ProfileService) *ProfileHandler {
	return &ProfileHandler{profile: profile}
}

func (h *ProfileHandler) GetProfile(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	profile, err := h.profile.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteAPIError(w, http.StatusNotFound, "not_found", "profile not found")
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to load profile")
		return
	}

	WriteJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) PatchMetrics(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.PatchMetricsRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	profile, err := h.profile.PatchMetrics(r.Context(), userID, req)
	if err != nil {
		WriteAPIError(w, http.StatusBadRequest, "validation_error", "failed to update metrics")
		return
	}
	WriteJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) PatchPreferences(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.PatchPreferencesRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	profile, err := h.profile.PatchPreferences(r.Context(), userID, req)
	if err != nil {
		WriteAPIError(w, http.StatusBadRequest, "validation_error", "failed to update preferences")
		return
	}
	WriteJSON(w, http.StatusOK, profile)
}

func (h *ProfileHandler) PatchBudget(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.PatchBudgetRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	profile, err := h.profile.PatchBudget(r.Context(), userID, req)
	if err != nil {
		if errors.Is(err, services.ErrInvalidBudget) {
			WriteAPIError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		WriteAPIError(w, http.StatusBadRequest, "validation_error", "failed to update budget")
		return
	}
	WriteJSON(w, http.StatusOK, profile)
}
