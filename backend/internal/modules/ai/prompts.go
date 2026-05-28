package ai

import (
	"errors"
	"fmt"
	"strings"
)

type PromptTemplate string

const (
	PromptTemplateSingleMeal PromptTemplate = "single_meal"
	PromptTemplateDayPlan    PromptTemplate = "day_plan"
	PromptTemplateWeekPlan   PromptTemplate = "week_plan"
	PromptTemplateMonthPlan  PromptTemplate = "month_plan"
)

var ErrUnknownPromptTemplate = errors.New("unknown prompt template")

type PromptRequest struct {
	UserRequest            string
	ResponseContract       string
	StartDate              string
	EndDate                string
	BudgetTarget           string
	BodyMetricsSummary     string
	PreferenceSummary      string
	PantrySnapshot         []string
	AdditionalInstructions []string
}

func BuildPrompt(template PromptTemplate, req PromptRequest) (string, error) {
	var objective string
	switch template {
	case PromptTemplateSingleMeal:
		objective = "Create a single meal recommendation."
	case PromptTemplateDayPlan:
		objective = "Create a 1-day meal plan with exactly 4 sections: breakfast, lunch, dinner, snacks."
	case PromptTemplateWeekPlan:
		objective = "Create a 7-day meal plan with exactly 4 sections per day: breakfast, lunch, dinner, snacks."
	case PromptTemplateMonthPlan:
		objective = "Create a month meal plan with exactly 4 sections per day: breakfast, lunch, dinner, snacks."
	default:
		return "", ErrUnknownPromptTemplate
	}

	sections := []string{
		"You are PantryPal's meal planning model.",
		objective,
		"Return valid JSON only.",
		"Do not include markdown fences, commentary, or explanatory text.",
		"Use the exact meal_type values breakfast, lunch, dinner, and snacks.",
	}

	if contract := strings.TrimSpace(req.ResponseContract); contract != "" {
		sections = append(sections, "Response contract:\n"+contract)
	}
	if request := strings.TrimSpace(req.UserRequest); request != "" {
		sections = append(sections, "User request:\n"+request)
	}
	if req.StartDate != "" || req.EndDate != "" {
		sections = append(sections, fmt.Sprintf("Plan window: start=%s end=%s", valueOrPlaceholder(req.StartDate), valueOrPlaceholder(req.EndDate)))
	}
	if budget := strings.TrimSpace(req.BudgetTarget); budget != "" {
		sections = append(sections, "Budget target:\n"+budget)
	}
	if metrics := strings.TrimSpace(req.BodyMetricsSummary); metrics != "" {
		sections = append(sections, "Body metrics:\n"+metrics)
	}
	if preferences := strings.TrimSpace(req.PreferenceSummary); preferences != "" {
		sections = append(sections, "Preferences and restrictions:\n"+preferences)
	}
	if len(req.PantrySnapshot) > 0 {
		sections = append(sections, "Pantry snapshot:\n- "+strings.Join(req.PantrySnapshot, "\n- "))
	}

	guardrails := []string{
		"Use realistic grocery ingredients and household units.",
		"Keep ingredient names specific and easy to match to USDA foods.",
		"Keep meals budget-aware and avoid luxury ingredients unless explicitly requested.",
	}
	if len(req.AdditionalInstructions) > 0 {
		guardrails = append(guardrails, req.AdditionalInstructions...)
	}
	sections = append(sections, "Guardrails:\n- "+strings.Join(guardrails, "\n- "))

	return strings.Join(sections, "\n\n"), nil
}

func valueOrPlaceholder(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return "n/a"
	}
	return value
}
