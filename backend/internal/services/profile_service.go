package services

import (
	"context"
	"errors"

	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

var ErrInvalidBudget = errors.New("amountCents is required and must be >= 0")

type ProfileService struct {
	users *repositories.UserRepository
}

func NewProfileService(users *repositories.UserRepository) *ProfileService {
	return &ProfileService{users: users}
}

func (s *ProfileService) GetProfile(ctx context.Context, userID string) (dto.ProfileResponse, error) {
	return s.users.GetProfile(ctx, userID)
}

func (s *ProfileService) PatchMetrics(ctx context.Context, userID string, req dto.PatchMetricsRequest) (dto.ProfileResponse, error) {
	if err := s.users.UpsertMetrics(ctx, userID, req); err != nil {
		return dto.ProfileResponse{}, err
	}
	return s.users.GetProfile(ctx, userID)
}

func (s *ProfileService) PatchPreferences(ctx context.Context, userID string, req dto.PatchPreferencesRequest) (dto.ProfileResponse, error) {
	if err := s.users.UpsertPreferences(ctx, userID, req); err != nil {
		return dto.ProfileResponse{}, err
	}
	return s.users.GetProfile(ctx, userID)
}

func (s *ProfileService) PatchBudget(ctx context.Context, userID string, req dto.PatchBudgetRequest) (dto.ProfileResponse, error) {
	if req.AmountCents == nil || *req.AmountCents < 0 {
		return dto.ProfileResponse{}, ErrInvalidBudget
	}
	if err := s.users.UpsertBudget(ctx, userID, req); err != nil {
		return dto.ProfileResponse{}, err
	}
	return s.users.GetProfile(ctx, userID)
}
