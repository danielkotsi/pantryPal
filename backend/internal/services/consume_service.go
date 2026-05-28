package services

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

var (
	ErrPlanMealNotFound    = errors.New("plan meal not found")
	ErrMealAlreadyConsumed = errors.New("meal has already been consumed")
)

type ConsumeService struct {
	plans    *repositories.PlanRepository
	recipes  *repositories.RecipeRepository
	pantry   *repositories.PantryRepository
	consumption *repositories.ConsumptionLogRepository
}

func NewConsumeService(plans *repositories.PlanRepository, recipes *repositories.RecipeRepository, pantry *repositories.PantryRepository, consumption *repositories.ConsumptionLogRepository) *ConsumeService {
	return &ConsumeService{plans: plans, recipes: recipes, pantry: pantry, consumption: consumption}
}

func (s *ConsumeService) ConsumeMeal(ctx context.Context, userID, mealID string) (dto.ConsumeMealResponse, error) {
	meal, err := s.plans.GetPlanMealByID(ctx, userID, mealID)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.ConsumeMealResponse{}, ErrPlanMealNotFound
	}
	if err != nil {
		return dto.ConsumeMealResponse{}, err
	}

	if meal.IsConsumed {
		return dto.ConsumeMealResponse{}, ErrMealAlreadyConsumed
	}

	var items []dto.ConsumedItemInfo

	if meal.RecipeID != "" {
		recipe, err := s.recipes.GetByID(ctx, meal.RecipeID)
		if err != nil {
			return dto.ConsumeMealResponse{}, fmt.Errorf("load recipe: %w", err)
		}

		scale := meal.Servings / float64(recipe.Servings)

		for _, ing := range recipe.Ingredients {
			deductQty := ing.Quantity * scale

			info := s.deductIngredient(ctx, userID, mealID, ing.FDCID, ing.Description, deductQty, ing.Unit)
			items = append(items, info)
		}
	}

	if err := s.plans.MarkMealConsumed(ctx, mealID); err != nil {
		return dto.ConsumeMealResponse{}, fmt.Errorf("mark meal consumed: %w", err)
	}

	now := time.Now().UTC().Format(time.RFC3339)

	return dto.ConsumeMealResponse{
		MealID:    mealID,
		Consumed:  true,
		ConsumedAt: now,
		Items:     items,
	}, nil
}

func (s *ConsumeService) deductIngredient(ctx context.Context, userID, mealID string, fdcID int64, description string, quantity float64, unit string) dto.ConsumedItemInfo {
	info := dto.ConsumedItemInfo{
		FDCID:           fdcID,
		Description:     description,
		QuantityDeducted: quantity,
		Unit:            unit,
	}

	deduction, err := s.pantry.DeductByFood(ctx, userID, fdcID, unit, quantity)
	if errors.Is(err, sql.ErrNoRows) {
		info.BeforeQuantity = 0
		info.AfterQuantity = 0
		info.QuantityDeducted = 0
		info.Warning = fmt.Sprintf("no pantry item found for %s (%s)", description, unit)

		_ = s.consumption.Insert(ctx, repositories.ConsumptionLogEntry{
			PlanMealID:       mealID,
			UserID:           userID,
			FDCID:            fdcID,
			QuantityDeducted: 0,
			Unit:             unit,
			Warning:          info.Warning,
		})
		return info
	}
	if err != nil {
		info.Warning = fmt.Sprintf("failed to deduct: %v", err)
		return info
	}

	info.PantryItemID = deduction.PantryItemID
	info.BeforeQuantity = deduction.BeforeQty
	info.AfterQuantity = deduction.AfterQty

	if deduction.BeforeQty < quantity {
		info.Warning = fmt.Sprintf("insufficient quantity: had %.2f %s, needed %.2f %s", deduction.BeforeQty, unit, quantity, unit)
		info.QuantityDeducted = deduction.BeforeQty - deduction.AfterQty
	}

	_ = s.consumption.Insert(ctx, repositories.ConsumptionLogEntry{
		PlanMealID:       mealID,
		UserID:           userID,
		FDCID:            fdcID,
		PantryItemID:     deduction.PantryItemID,
		QuantityDeducted: info.QuantityDeducted,
		Unit:             unit,
		BeforeQuantity:   deduction.BeforeQty,
		AfterQuantity:    deduction.AfterQty,
		Warning:          info.Warning,
	})

	return info
}
