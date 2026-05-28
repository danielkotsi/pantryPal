package handlers

import (
	"errors"
	"net/http"

	"pantrypal/backend/internal/services"
	"pantrypal/backend/internal/transport/http/dto"
	"pantrypal/backend/internal/transport/http/middleware"
)

type PlanHandler struct {
	plans *services.PlanService
}

func NewPlanHandler(plans *services.PlanService) *PlanHandler {
	return &PlanHandler{plans: plans}
}

func (h *PlanHandler) CreateProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.PlanProposalRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	res, err := h.plans.CreateProposal(r.Context(), userID, req)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPlanType),
			errors.Is(err, services.ErrInvalidStartDate),
			errors.Is(err, services.ErrInvalidEndDate),
			errors.Is(err, services.ErrInvalidPlanMeals),
			errors.Is(err, services.ErrInvalidMealSection):
			WriteAPIError(w, http.StatusBadRequest, "validation_error", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to create plan proposal")
		}
		return
	}

	WriteJSON(w, http.StatusCreated, res)
}

func (h *PlanHandler) AcceptProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	res, err := h.plans.AcceptProposal(r.Context(), userID, r.PathValue("id"))
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPlanStatus):
			WriteAPIError(w, http.StatusNotFound, "not_found", err.Error())
		case errors.Is(err, services.ErrPlanNotProposal):
			WriteAPIError(w, http.StatusConflict, "conflict", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to accept plan")
		}
		return
	}

	WriteJSON(w, http.StatusOK, res)
}

func (h *PlanHandler) DeclineProposal(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	var req dto.DeclinePlanRequest
	if err := ReadJSON(r, &req); err != nil {
		WriteAPIError(w, http.StatusBadRequest, "bad_request", "invalid JSON payload")
		return
	}

	res, err := h.plans.DeclineProposal(r.Context(), userID, r.PathValue("id"), req.Reason)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrInvalidPlanStatus):
			WriteAPIError(w, http.StatusNotFound, "not_found", err.Error())
		case errors.Is(err, services.ErrPlanNotProposal):
			WriteAPIError(w, http.StatusConflict, "conflict", err.Error())
		default:
			WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to decline plan")
		}
		return
	}

	WriteJSON(w, http.StatusOK, res)
}

func (h *PlanHandler) GetWeekPlan(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.UserIDFromContext(r.Context())
	if !ok {
		WriteAPIError(w, http.StatusUnauthorized, "unauthorized", "missing auth context")
		return
	}

	res, err := h.plans.GetWeekPlan(r.Context(), userID, r.URL.Query().Get("start"))
	if err != nil {
		if errors.Is(err, services.ErrWeekStartDate) {
			WriteAPIError(w, http.StatusBadRequest, "validation_error", err.Error())
			return
		}
		WriteAPIError(w, http.StatusInternalServerError, "internal_error", "failed to load week plan")
		return
	}

	WriteJSON(w, http.StatusOK, res)
}
