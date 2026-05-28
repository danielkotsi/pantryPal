package handlers

import (
	"net/http"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/dto"
	"pantrypal/backend/internal/transport/http/middleware"
)

type GenerateHandler struct {
	generate *services.GenerateService
}

func NewGenerateHandler(generate *services.GenerateService) *GenerateHandler {
	return &GenerateHandler{generate: generate}
}

func (h *GenerateHandler) GeneratePlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.GeneratePlanRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}
	if req.PeriodType == "" {
		req.PeriodType = "week"
	}

	result, err := h.generate.GeneratePlan(r.Context(), userID, req.PeriodType, req.Message)
	if err != nil {
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to generate plan")
		return
	}

	WriteJSON(w, http.StatusCreated, dto.GeneratePlanResponse{
		Proposal:       result.Proposal,
		FallbackActive: result.FallbackActive,
	})
}
