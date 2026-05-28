package handlers

import (
	"database/sql"
	"errors"
	"net/http"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/dto"
	"pantrypal/backend/internal/transport/http/middleware"
)

type AuthHandler struct {
	auth *services.AuthService
}

func NewAuthHandler(auth *services.AuthService) *AuthHandler {
	return &AuthHandler{auth: auth}
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.RegisterRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	res, err := h.auth.Register(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidEmail), errors.Is(err, services.ErrWeakPassword):
			WriteAPIError(w, http.StatusBadRequest, "validation_error", err.Error())
		case errors.Is(err, services.ErrEmailConflict):
			WriteAPIError(w, http.StatusConflict, "conflict", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to register user")
		}
		return
	}

	WriteJSON(w, http.StatusCreated, res)
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	res, err := h.auth.Login(r.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidCredentials):
			WriteAPIError(w, http.StatusUnauthorized, "unauthorized", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to login")
		}
		return
	}

	WriteJSON(w, http.StatusOK, res)
}

func (h *AuthHandler) Me(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	me, err := h.auth.Me(r.Context(), userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "user not found")
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to load current user")
		return
	}

	WriteJSON(w, http.StatusOK, map[string]dto.UserResponse{"user": me})
}
