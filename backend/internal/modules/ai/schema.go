package ai

import "pantrypal/backend/internal/transport/http/dto"

type PlanResponse struct {
	PeriodType      string            `json:"period_type"`
	StartDate       string            `json:"start_date"`
	EndDate         string            `json:"end_date,omitempty"`
	Currency        string            `json:"currency,omitempty"`
	EstimatedTotal  float64           `json:"estimated_total_cost"`
	Days            []DayPlan         `json:"days"`
	Notes           string            `json:"notes,omitempty"`
	RequestMetadata map[string]string `json:"request_metadata,omitempty"`
}

type DayPlan struct {
	Date  string       `json:"date"`
	Meals []MealDetail `json:"meals"`
}

type MealDetail struct {
	MealType      string            `json:"meal_type"`
	RecipeName    string            `json:"recipe_name"`
	Servings      float64           `json:"servings"`
	Ingredients   []IngredientInput `json:"ingredients"`
	Macros        MacroInput        `json:"macros"`
	EstimatedCost float64           `json:"estimated_cost"`
	Notes         string            `json:"notes,omitempty"`
}

type IngredientInput struct {
	Name     string  `json:"name"`
	Quantity float64 `json:"quantity"`
	Unit     string  `json:"unit"`
}

type MacroInput struct {
	Protein  float64 `json:"protein"`
	Carbs    float64 `json:"carbs"`
	Fat      float64 `json:"fat"`
	Calories float64 `json:"calories"`
}

type NormalizedPlan struct {
	Proposal         dto.PlanProposalRequest
	Ingredients      []NormalizedIngredient
	UnmatchedReasons []string
	Parse            ParseMetadata
}

type NormalizedIngredient struct {
	ScheduledDate string
	MealSection   string
	RecipeName    string
	Name          string
	Quantity      float64
	Unit          string
}

type ParseMetadata struct {
	UsedRepair     bool
	RepairActions  []string
	ValidationStep string
}
