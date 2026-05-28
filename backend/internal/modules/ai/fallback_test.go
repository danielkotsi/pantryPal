package ai

import (
	"strings"
	"testing"
)

func TestFallbackWeekPlan(t *testing.T) {
	parsed, err := FallbackWeekPlan("2026-06-01")
	if err != nil {
		t.Fatalf("FallbackWeekPlan returned error: %v", err)
	}

	if parsed.Proposal.Source != "fallback" {
		t.Fatalf("expected source fallback, got %q", parsed.Proposal.Source)
	}
	if parsed.Proposal.PeriodType != "week" {
		t.Fatalf("expected week period type, got %q", parsed.Proposal.PeriodType)
	}
	if parsed.Proposal.StartDate != "2026-06-01" {
		t.Fatalf("expected start date 2026-06-01, got %q", parsed.Proposal.StartDate)
	}
	if parsed.Proposal.EndDate != "2026-06-07" {
		t.Fatalf("expected end date 2026-06-07, got %q", parsed.Proposal.EndDate)
	}

	if len(parsed.Proposal.Meals) != 28 {
		t.Fatalf("expected 28 meals (7 days x 4 sections), got %d", len(parsed.Proposal.Meals))
	}

	if parsed.Proposal.AICostCentsTotal <= 0 {
		t.Fatalf("expected positive total cost cents, got %d", parsed.Proposal.AICostCentsTotal)
	}

	mealDates := make(map[string]int)
	mealSections := make(map[string]int)
	for _, meal := range parsed.Proposal.Meals {
		mealDates[meal.ScheduledDate]++
		mealSections[meal.MealSection]++
	}

	if len(mealDates) != 7 {
		t.Fatalf("expected meals spread across 7 unique dates, got %d dates", len(mealDates))
	}
	for date, count := range mealDates {
		if count != 4 {
			t.Fatalf("expected exactly 4 meals on date %q, got %d", date, count)
		}
	}

	expectedSections := []string{"breakfast", "lunch", "dinner", "snacks"}
	for _, section := range expectedSections {
		if mealSections[section] != 7 {
			t.Fatalf("expected 7 %q meals, got %d", section, mealSections[section])
		}
	}

	if len(parsed.Ingredients) == 0 {
		t.Fatal("expected non-empty ingredients list")
	}

	for _, meal := range parsed.Proposal.Meals {
		if meal.RecipeName == "" {
			t.Fatal("all meals must have a recipe_name")
		}
		if meal.Servings <= 0 {
			t.Fatal("all meals must have positive servings")
		}
		if meal.Macros.Calories <= 0 {
			t.Fatal("all meals must have positive calories")
		}
	}

	if parsed.Parse.UsedRepair {
		t.Fatalf("fallback output should not trigger repairs, got actions: %s",
			strings.Join(parsed.Parse.RepairActions, ", "))
	}
}

func TestFallbackWeekPlanDefaultsDate(t *testing.T) {
	parsed, err := FallbackWeekPlan("")
	if err != nil {
		t.Fatalf("FallbackWeekPlan with empty date returned error: %v", err)
	}

	if len(parsed.Proposal.Meals) != 28 {
		t.Fatalf("expected 28 meals even with default date, got %d", len(parsed.Proposal.Meals))
	}
}

func TestFallbackWeekPlanRotation(t *testing.T) {
	parsed1, err := FallbackWeekPlan("2026-06-01")
	if err != nil {
		t.Fatalf("first call error: %v", err)
	}

	parsed2, err := FallbackWeekPlan("2026-06-08")
	if err != nil {
		t.Fatalf("second call error: %v", err)
	}

	if parsed1.Proposal.StartDate == parsed2.Proposal.StartDate {
		t.Fatal("two different start dates should produce different plans")
	}

	if parsed1.Proposal.AICostCentsTotal <= 0 || parsed2.Proposal.AICostCentsTotal <= 0 {
		t.Fatal("both plans should have positive total cost")
	}
}
