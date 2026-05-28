package ai

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"pantrypal/backend/internal/transport/http/dto"
)

var (
	ErrInvalidAIJSON        = errors.New("ai response is not valid JSON")
	ErrInvalidAISchema      = errors.New("ai response failed schema validation")
	ErrUnsupportedPlanType  = errors.New("period_type must be one of meal, day, week, month")
	ErrInvalidMealType      = errors.New("meal_type must be one of breakfast, lunch, dinner, snacks (the alias snack is also accepted)")
	ErrMissingDays          = errors.New("days must not be empty")
	ErrMissingMeals         = errors.New("day meals must not be empty")
	ErrInvalidDate          = errors.New("dates must use YYYY-MM-DD")
	ErrInvalidServings      = errors.New("servings must be greater than zero")
	ErrInvalidIngredient    = errors.New("ingredients must include name, quantity, and unit")
	ErrInvalidMacros        = errors.New("macros must include non-negative protein, carbs, fat, and calories")
	ErrInvalidEstimatedCost = errors.New("estimated_cost must be zero or positive")
)

type ParseError struct {
	Reason string
	Field  string
}

func (e ParseError) Error() string {
	if e.Field == "" {
		return e.Reason
	}
	return fmt.Sprintf("%s: %s", e.Field, e.Reason)
}

func (e ParseError) Unwrap() error {
	return ErrInvalidAISchema
}

func ParsePlanResponse(raw string) (NormalizedPlan, error) {
	cleaned, repairs := repairRawJSON(raw)

	var payload PlanResponse
	if err := decodeStrict(cleaned, &payload); err != nil {
		fallback, convertRepairs, convErr := decodeWithNumericRepair(cleaned)
		if convErr != nil {
			return NormalizedPlan{}, fmt.Errorf("%w: %v", ErrInvalidAIJSON, err)
		}
		payload = fallback
		repairs = append(repairs, convertRepairs...)
	}

	normalized, err := normalizePlan(payload)
	if err != nil {
		return NormalizedPlan{}, err
	}
	normalized.Parse = ParseMetadata{
		UsedRepair:     len(repairs) > 0,
		RepairActions:  repairs,
		ValidationStep: "normalized",
	}
	return normalized, nil
}

func decodeStrict(raw string, dest *PlanResponse) error {
	dec := json.NewDecoder(strings.NewReader(raw))
	dec.DisallowUnknownFields()
	return dec.Decode(dest)
}

func decodeWithNumericRepair(raw string) (PlanResponse, []string, error) {
	var generic any
	if err := json.Unmarshal([]byte(raw), &generic); err != nil {
		return PlanResponse{}, nil, err
	}
	repaired, actions := coerceNumericStrings(generic, "")
	bytes, err := json.Marshal(repaired)
	if err != nil {
		return PlanResponse{}, nil, err
	}
	var payload PlanResponse
	if err := decodeStrict(string(bytes), &payload); err != nil {
		return PlanResponse{}, nil, err
	}
	return payload, actions, nil
}

func normalizePlan(payload PlanResponse) (NormalizedPlan, error) {
	periodType := strings.TrimSpace(strings.ToLower(payload.PeriodType))
	if periodType == "" {
		periodType = inferPeriodType(payload)
	}
	switch periodType {
	case "meal", "day", "week", "month":
	default:
		return NormalizedPlan{}, ParseError{Field: "period_type", Reason: ErrUnsupportedPlanType.Error()}
	}

	if len(payload.Days) == 0 {
		return NormalizedPlan{}, ParseError{Field: "days", Reason: ErrMissingDays.Error()}
	}

	notes := strings.TrimSpace(payload.Notes)
	proposal := dto.PlanProposalRequest{
		PeriodType:      planProposalPeriodType(periodType),
		StartDate:       strings.TrimSpace(payload.StartDate),
		EndDate:         strings.TrimSpace(payload.EndDate),
		Source:          "ai",
		ProposalVersion: 1,
		Notes:           notes,
	}

	ingredients := make([]NormalizedIngredient, 0)
	totalCostCents := 0
	seenSlots := make(map[string]struct{})
	sortedDays := append([]DayPlan(nil), payload.Days...)
	sort.Slice(sortedDays, func(i, j int) bool {
		return sortedDays[i].Date < sortedDays[j].Date
	})

	for dayIndex, day := range sortedDays {
		date := strings.TrimSpace(day.Date)
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return NormalizedPlan{}, ParseError{Field: fmt.Sprintf("days[%d].date", dayIndex), Reason: ErrInvalidDate.Error()}
		}
		if len(day.Meals) == 0 {
			return NormalizedPlan{}, ParseError{Field: fmt.Sprintf("days[%d].meals", dayIndex), Reason: ErrMissingMeals.Error()}
		}
		for mealIndex, meal := range day.Meals {
			mealField := fmt.Sprintf("days[%d].meals[%d]", dayIndex, mealIndex)
			mealSection, err := normalizeMealType(meal.MealType)
			if err != nil {
				return NormalizedPlan{}, ParseError{Field: mealField + ".meal_type", Reason: err.Error()}
			}
			if strings.TrimSpace(meal.RecipeName) == "" {
				return NormalizedPlan{}, ParseError{Field: mealField + ".recipe_name", Reason: "recipe_name is required"}
			}
			if meal.Servings <= 0 {
				return NormalizedPlan{}, ParseError{Field: mealField + ".servings", Reason: ErrInvalidServings.Error()}
			}
			if err := validateMacros(mealField+".macros", meal.Macros); err != nil {
				return NormalizedPlan{}, err
			}
			if meal.EstimatedCost < 0 {
				return NormalizedPlan{}, ParseError{Field: mealField + ".estimated_cost", Reason: ErrInvalidEstimatedCost.Error()}
			}
			if len(meal.Ingredients) == 0 {
				return NormalizedPlan{}, ParseError{Field: mealField + ".ingredients", Reason: "ingredients must not be empty"}
			}

			slotKey := date + "|" + mealSection
			if _, exists := seenSlots[slotKey]; exists {
				return NormalizedPlan{}, ParseError{Field: mealField + ".meal_type", Reason: "duplicate meal section for date"}
			}
			seenSlots[slotKey] = struct{}{}

			for ingredientIndex, ingredient := range meal.Ingredients {
				name := strings.TrimSpace(ingredient.Name)
				unit := strings.TrimSpace(ingredient.Unit)
				if name == "" || unit == "" || ingredient.Quantity <= 0 {
					return NormalizedPlan{}, ParseError{Field: fmt.Sprintf("%s.ingredients[%d]", mealField, ingredientIndex), Reason: ErrInvalidIngredient.Error()}
				}
				ingredients = append(ingredients, NormalizedIngredient{
					ScheduledDate: date,
					MealSection:   mealSection,
					RecipeName:    strings.TrimSpace(meal.RecipeName),
					Name:          name,
					Quantity:      ingredient.Quantity,
					Unit:          unit,
				})
			}

			costCents := dollarsToCents(meal.EstimatedCost)
			totalCostCents += costCents
			proposal.Meals = append(proposal.Meals, dto.PlanMealInput{
				ScheduledDate:      date,
				MealSection:        mealSection,
				RecipeName:         strings.TrimSpace(meal.RecipeName),
				Servings:           meal.Servings,
				EstimatedCostCents: costCents,
				Macros: dto.RecipeMacrosResponse{
					Calories: meal.Macros.Calories,
					ProteinG: meal.Macros.Protein,
					CarbsG:   meal.Macros.Carbs,
					FatG:     meal.Macros.Fat,
				},
			})
		}
	}

	if proposal.StartDate == "" {
		proposal.StartDate = sortedDays[0].Date
	}
	if proposal.EndDate == "" {
		proposal.EndDate = sortedDays[len(sortedDays)-1].Date
	}
	if _, err := time.Parse("2006-01-02", proposal.StartDate); err != nil {
		return NormalizedPlan{}, ParseError{Field: "start_date", Reason: ErrInvalidDate.Error()}
	}
	if _, err := time.Parse("2006-01-02", proposal.EndDate); err != nil {
		return NormalizedPlan{}, ParseError{Field: "end_date", Reason: ErrInvalidDate.Error()}
	}
	if proposal.EndDate < proposal.StartDate {
		return NormalizedPlan{}, ParseError{Field: "end_date", Reason: "end_date must not be earlier than start_date"}
	}

	proposal.AICostCentsTotal = totalCostCents
	return NormalizedPlan{Proposal: proposal, Ingredients: ingredients}, nil
}

func validateMacros(field string, macros MacroInput) error {
	if macros.Protein < 0 || macros.Carbs < 0 || macros.Fat < 0 || macros.Calories < 0 {
		return ParseError{Field: field, Reason: ErrInvalidMacros.Error()}
	}
	return nil
}

func normalizeMealType(raw string) (string, error) {
	mealType := strings.TrimSpace(strings.ToLower(raw))
	switch mealType {
	case "breakfast", "lunch", "dinner":
		return mealType, nil
	case "snack", "snacks":
		// Keep one canonical downstream value while tolerating the singular alias from model output.
		return "snacks", nil
	default:
		return "", ErrInvalidMealType
	}
}

func planProposalPeriodType(periodType string) string {
	if periodType == "meal" {
		return "day"
	}
	return periodType
}

func inferPeriodType(payload PlanResponse) string {
	if len(payload.Days) == 1 && len(payload.Days[0].Meals) == 1 {
		return "meal"
	}
	if len(payload.Days) == 1 {
		return "day"
	}
	if len(payload.Days) <= 7 {
		return "week"
	}
	return "month"
}

func dollarsToCents(value float64) int {
	return int(math.Round(value * 100))
}

func repairRawJSON(raw string) (string, []string) {
	trimmed := strings.TrimSpace(raw)
	actions := make([]string, 0)
	if strings.HasPrefix(trimmed, "```") {
		trimmed = strings.TrimPrefix(trimmed, "```json")
		trimmed = strings.TrimPrefix(trimmed, "```")
		trimmed = strings.TrimSuffix(strings.TrimSpace(trimmed), "```")
		actions = append(actions, "stripped_markdown_fences")
	}
	return strings.TrimSpace(trimmed), actions
}

func coerceNumericStrings(value any, path string) (any, []string) {
	actions := make([]string, 0)
	switch typed := value.(type) {
	case map[string]any:
		result := make(map[string]any, len(typed))
		for key, item := range typed {
			nextPath := key
			if path != "" {
				nextPath = path + "." + key
			}
			converted, nestedActions := coerceNumericStrings(item, nextPath)
			result[key] = converted
			actions = append(actions, nestedActions...)
		}
		return result, actions
	case []any:
		result := make([]any, len(typed))
		for index, item := range typed {
			nextPath := fmt.Sprintf("%s[%d]", path, index)
			converted, nestedActions := coerceNumericStrings(item, nextPath)
			result[index] = converted
			actions = append(actions, nestedActions...)
		}
		return result, actions
	case string:
		trimmed := strings.TrimSpace(typed)
		if trimmed == "" {
			return typed, nil
		}
		if looksNumericPath(path) {
			if n, err := strconv.ParseFloat(trimmed, 64); err == nil {
				actions = append(actions, "coerced_numeric_string:"+path)
				return n, actions
			}
		}
		return typed, nil
	default:
		return value, nil
	}
}

func looksNumericPath(path string) bool {
	for _, suffix := range []string{
		"servings",
		"estimated_cost",
		"estimated_total_cost",
		"quantity",
		"protein",
		"carbs",
		"fat",
		"calories",
	} {
		if strings.HasSuffix(path, suffix) {
			return true
		}
	}
	return false
}
