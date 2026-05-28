package services

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

var (
	ErrInvalidSearchQuery = errors.New("query parameter q is required")
	ErrInvalidPantryFood  = errors.New("fdcId must reference an existing USDA food")
	ErrInvalidQuantity    = errors.New("quantity must be greater than 0")
	ErrInvalidUnit        = errors.New("unit is required")
	ErrPantryItemNotFound = errors.New("pantry item not found")
)

type PantryService struct {
	foods  *repositories.FoodRepository
	pantry *repositories.PantryRepository
}

func NewPantryService(foods *repositories.FoodRepository, pantry *repositories.PantryRepository) *PantryService {
	return &PantryService{foods: foods, pantry: pantry}
}

func (s *PantryService) SearchFoods(ctx context.Context, query string) ([]dto.FoodSearchItem, error) {
	if strings.TrimSpace(query) == "" {
		return nil, ErrInvalidSearchQuery
	}
	return s.foods.SearchFoods(ctx, query, 20)
}

func (s *PantryService) ListPantryItems(ctx context.Context, userID string) ([]dto.PantryItemResponse, error) {
	return s.pantry.ListByUserID(ctx, userID)
}

func (s *PantryService) AddPantryItem(ctx context.Context, userID string, req dto.PantryItemRequest) (dto.PantryItemResponse, error) {
	if req.Quantity <= 0 {
		return dto.PantryItemResponse{}, ErrInvalidQuantity
	}
	if strings.TrimSpace(req.Unit) == "" {
		return dto.PantryItemResponse{}, ErrInvalidUnit
	}
	exists, err := s.foods.ExistsByFDCID(ctx, req.FDCID)
	if err != nil {
		return dto.PantryItemResponse{}, err
	}
	if !exists {
		return dto.PantryItemResponse{}, ErrInvalidPantryFood
	}
	return s.pantry.Upsert(ctx, userID, req)
}

func (s *PantryService) PatchPantryItem(ctx context.Context, userID, itemID string, req dto.PantryItemPatchRequest) (dto.PantryItemResponse, error) {
	item, err := s.pantry.PatchQuantity(ctx, userID, itemID, req.QuantityDelta)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.PantryItemResponse{}, ErrPantryItemNotFound
	}
	return item, err
}

func (s *PantryService) DeletePantryItem(ctx context.Context, userID, itemID string) error {
	err := s.pantry.Delete(ctx, userID, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return ErrPantryItemNotFound
	}
	return err
}
