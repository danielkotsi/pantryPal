package ai

import (
	"errors"
	"strings"
	"testing"
)

func TestParsePlanResponseValidWeekPlan(t *testing.T) {
	raw := `{
		"period_type": "week",
		"start_date": "2026-06-01",
		"end_date": "2026-06-07",
		"estimated_total_cost": 48.75,
		"days": [
			{
				"date": "2026-06-01",
				"meals": [
					{
						"meal_type": "breakfast",
						"recipe_name": "Greek Yogurt Bowl",
						"servings": 1,
						"ingredients": [
							{"name": "Yogurt, Greek, plain, nonfat", "quantity": 1, "unit": "cup"}
						],
						"macros": {"protein": 20, "carbs": 18, "fat": 0, "calories": 180},
						"estimated_cost": 2.5
					},
					{
						"meal_type": "snack",
						"recipe_name": "Apple and Peanut Butter",
						"servings": 1,
						"ingredients": [
							{"name": "Apple", "quantity": 1, "unit": "piece"},
							{"name": "Peanut butter", "quantity": 2, "unit": "tbsp"}
						],
						"macros": {"protein": 8, "carbs": 28, "fat": 16, "calories": 280},
						"estimated_cost": 1.75
					}
				]
			}
		]
	}`

	parsed, err := ParsePlanResponse(raw)
	if err != nil {
		t.Fatalf("ParsePlanResponse returned error: %v", err)
	}

	if parsed.Proposal.PeriodType != "week" {
		t.Fatalf("expected week period type, got %q", parsed.Proposal.PeriodType)
	}
	if len(parsed.Proposal.Meals) != 2 {
		t.Fatalf("expected 2 meals, got %d", len(parsed.Proposal.Meals))
	}
	if parsed.Proposal.Meals[1].MealSection != "snacks" {
		t.Fatalf("expected snack alias to normalize to canonical snacks, got %q", parsed.Proposal.Meals[1].MealSection)
	}
	if parsed.Proposal.AICostCentsTotal != 425 {
		t.Fatalf("expected total cost 425 cents, got %d", parsed.Proposal.AICostCentsTotal)
	}
}

func TestParsePlanResponseRepairsMarkdownAndNumericStrings(t *testing.T) {
	raw := "```json\n" + `{
		"period_type": "day",
		"start_date": "2026-06-01",
		"days": [
			{
				"date": "2026-06-01",
				"meals": [
					{
						"meal_type": "dinner",
						"recipe_name": "Chicken Rice Bowl",
						"servings": "2",
						"ingredients": [
							{"name": "Chicken breast", "quantity": "12", "unit": "oz"}
						],
						"macros": {"protein": "45", "carbs": "52", "fat": "8", "calories": "460"},
						"estimated_cost": "6.25"
					}
				]
			}
		]
	}` + "\n```"

	parsed, err := ParsePlanResponse(raw)
	if err != nil {
		t.Fatalf("ParsePlanResponse returned error: %v", err)
	}
	if !parsed.Parse.UsedRepair {
		t.Fatal("expected parser to record repair usage")
	}
	joined := strings.Join(parsed.Parse.RepairActions, ",")
	if !strings.Contains(joined, "stripped_markdown_fences") {
		t.Fatalf("expected markdown fence repair, got %q", joined)
	}
	if !strings.Contains(joined, "coerced_numeric_string:days[0].meals[0].servings") {
		t.Fatalf("expected numeric coercion repair, got %q", joined)
	}
	if parsed.Proposal.Meals[0].EstimatedCostCents != 625 {
		t.Fatalf("expected repaired cost cents 625, got %d", parsed.Proposal.Meals[0].EstimatedCostCents)
	}
}

func TestParsePlanResponseRejectsInvalidMealType(t *testing.T) {
	raw := `{
		"period_type": "day",
		"start_date": "2026-06-01",
		"days": [
			{
				"date": "2026-06-01",
				"meals": [
					{
						"meal_type": "dessert",
						"recipe_name": "Cake",
						"servings": 1,
						"ingredients": [
							{"name": "Flour", "quantity": 1, "unit": "cup"}
						],
						"macros": {"protein": 1, "carbs": 10, "fat": 5, "calories": 100},
						"estimated_cost": 1
					}
				]
			}
		]
	}`

	_, err := ParsePlanResponse(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidAISchema) {
		t.Fatalf("expected schema error, got %v", err)
	}
	if !strings.Contains(err.Error(), "meal_type") {
		t.Fatalf("expected field-specific error, got %v", err)
	}
}

func TestParsePlanResponseRejectsUnknownFields(t *testing.T) {
	raw := `{
		"period_type": "day",
		"start_date": "2026-06-01",
		"days": [
			{
				"date": "2026-06-01",
				"meals": [
					{
						"meal_type": "breakfast",
						"recipe_name": "Oats",
						"servings": 1,
						"ingredients": [
							{"name": "Oats", "quantity": 1, "unit": "cup"}
						],
						"macros": {"protein": 5, "carbs": 27, "fat": 3, "calories": 150},
						"estimated_cost": 0.9,
						"unexpected": true
					}
				]
			}
		]
	}`

	_, err := ParsePlanResponse(raw)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, ErrInvalidAIJSON) {
		t.Fatalf("expected invalid JSON/schema decode error, got %v", err)
	}
}
