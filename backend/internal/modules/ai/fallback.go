package ai

import (
	"encoding/json"
	"fmt"
	"time"
)

func FallbackWeekPlan(startDate string) (NormalizedPlan, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		start = time.Now().UTC().Truncate(24 * time.Hour)
	}
	end := start.AddDate(0, 0, 6)

	rotation := buildRotation()

	totalCost := 0.0
	days := make([]DayPlan, 7)
	for i := 0; i < 7; i++ {
		date := start.AddDate(0, 0, i).Format("2006-01-02")
		dayMeals := rotation[i%len(rotation)]
		meals := make([]MealDetail, 4)
		copy(meals, dayMeals)
		for mi := range meals {
			totalCost += meals[mi].EstimatedCost
		}
		days[i] = DayPlan{Date: date, Meals: meals}
	}

	payload := PlanResponse{
		PeriodType:     "week",
		StartDate:      start.Format("2006-01-02"),
		EndDate:        end.Format("2006-01-02"),
		EstimatedTotal: totalCost,
		Days:           days,
		Notes:          "Fallback plan generated from seeded recipes",
	}

	raw, err := json.Marshal(payload)
	if err != nil {
		return NormalizedPlan{}, fmt.Errorf("fallback marshal: %w", err)
	}

	normalized, err := ParsePlanResponse(string(raw))
	if err != nil {
		return NormalizedPlan{}, fmt.Errorf("fallback parse: %w", err)
	}

	normalized.Proposal.Source = "fallback"
	return normalized, nil
}

func buildRotation() [][]MealDetail {
	breakfastA := MealDetail{
		MealType:   "breakfast",
		RecipeName: "Banana Oat Bowl",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Oats", Quantity: 80, Unit: "g"},
			{Name: "Milk", Quantity: 200, Unit: "ml"},
			{Name: "Banana", Quantity: 120, Unit: "g"},
		},
		Macros: MacroInput{
			Protein:  17,
			Carbs:    69,
			Fat:      10,
			Calories: 430,
		},
		EstimatedCost: 1.80,
	}

	breakfastB := MealDetail{
		MealType:   "breakfast",
		RecipeName: "Yogurt and Apple Bowl",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Yogurt, Greek, plain, nonfat", Quantity: 200, Unit: "g"},
			{Name: "Apple", Quantity: 160, Unit: "g"},
		},
		Macros: MacroInput{
			Protein:  18,
			Carbs:    40,
			Fat:      1,
			Calories: 280,
		},
		EstimatedCost: 2.10,
	}

	lunchA := MealDetail{
		MealType:   "lunch",
		RecipeName: "Chicken Rice Plate",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Chicken breast", Quantity: 180, Unit: "g"},
			{Name: "Rice, white, cooked", Quantity: 170, Unit: "g"},
			{Name: "Broccoli, cooked", Quantity: 100, Unit: "g"},
			{Name: "Olive oil", Quantity: 10, Unit: "ml"},
		},
		Macros: MacroInput{
			Protein:  52,
			Carbs:    57,
			Fat:      20,
			Calories: 640,
		},
		EstimatedCost: 4.20,
	}

	lunchB := MealDetail{
		MealType:   "lunch",
		RecipeName: "Egg and Rice Bowl",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Egg, whole, cooked", Quantity: 2, Unit: "piece"},
			{Name: "Rice, white, cooked", Quantity: 150, Unit: "g"},
			{Name: "Broccoli, cooked", Quantity: 80, Unit: "g"},
			{Name: "Olive oil", Quantity: 8, Unit: "ml"},
		},
		Macros: MacroInput{
			Protein:  24,
			Carbs:    63,
			Fat:      22,
			Calories: 560,
		},
		EstimatedCost: 2.60,
	}

	dinnerA := MealDetail{
		MealType:   "dinner",
		RecipeName: "Egg Fried Rice Lite",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Egg, whole, cooked", Quantity: 2, Unit: "piece"},
			{Name: "Rice, white, cooked", Quantity: 150, Unit: "g"},
			{Name: "Broccoli, cooked", Quantity: 80, Unit: "g"},
			{Name: "Olive oil", Quantity: 8, Unit: "ml"},
		},
		Macros: MacroInput{
			Protein:  24,
			Carbs:    63,
			Fat:      22,
			Calories: 560,
		},
		EstimatedCost: 2.60,
	}

	dinnerB := MealDetail{
		MealType:   "dinner",
		RecipeName: "Grilled Chicken and Rice",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Chicken breast", Quantity: 180, Unit: "g"},
			{Name: "Rice, white, cooked", Quantity: 170, Unit: "g"},
			{Name: "Broccoli, cooked", Quantity: 100, Unit: "g"},
			{Name: "Olive oil", Quantity: 10, Unit: "ml"},
		},
		Macros: MacroInput{
			Protein:  52,
			Carbs:    57,
			Fat:      20,
			Calories: 640,
		},
		EstimatedCost: 4.20,
	}

	snackA := MealDetail{
		MealType:   "snacks",
		RecipeName: "Yogurt Apple Snack",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Yogurt, Greek, plain, nonfat", Quantity: 200, Unit: "g"},
			{Name: "Apple", Quantity: 160, Unit: "g"},
		},
		Macros: MacroInput{
			Protein:  16,
			Carbs:    30,
			Fat:      1,
			Calories: 210,
		},
		EstimatedCost: 1.70,
	}

	snackB := MealDetail{
		MealType:   "snacks",
		RecipeName: "Banana and Yogurt",
		Servings:   1,
		Ingredients: []IngredientInput{
			{Name: "Yogurt, Greek, plain, nonfat", Quantity: 150, Unit: "g"},
			{Name: "Banana", Quantity: 120, Unit: "g"},
		},
		Macros: MacroInput{
			Protein:  14,
			Carbs:    35,
			Fat:      1,
			Calories: 230,
		},
		EstimatedCost: 1.60,
	}

	dayOdd := []MealDetail{breakfastA, lunchA, dinnerA, snackA}
	dayEven := []MealDetail{breakfastB, lunchB, dinnerB, snackB}

	return [][]MealDetail{dayOdd, dayEven}
}
