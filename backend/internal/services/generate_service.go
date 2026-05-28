package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"pantrypal/backend/internal/modules/ai"
	"pantrypal/backend/internal/transport/http/dto"
)

type GeneratePlanResult struct {
	Proposal       dto.ProposalResponse
	FallbackActive bool
}

type GenerateService struct {
	client  *ai.Client
	plans   *PlanService
	profile *ProfileService
	pantry  *PantryService
}

func NewGenerateService(client *ai.Client, plans *PlanService, profile *ProfileService, pantry *PantryService) *GenerateService {
	return &GenerateService{client: client, plans: plans, profile: profile, pantry: pantry}
}

func (s *GenerateService) GeneratePlan(ctx context.Context, userID string, periodType string, userMessage string) (GeneratePlanResult, error) {
	if s.client != nil {
		result, err := s.tryGemini(ctx, userID, periodType, userMessage)
		if err == nil {
			return result, nil
		}
	}

	return s.fallback(ctx, userID)
}

func (s *GenerateService) tryGemini(ctx context.Context, userID string, periodType string, userMessage string) (GeneratePlanResult, error) {
	promptReq := buildPromptContext(ctx, s.profile, s.pantry, userID, periodType, userMessage)
	promptReq.ResponseContract = responseContract(periodType)

	prompt, err := ai.BuildPrompt(aiPromptTemplate(periodType), promptReq)
	if err != nil {
		return GeneratePlanResult{}, fmt.Errorf("build prompt: %w", err)
	}

	geminiResp, err := s.client.Generate(ctx, ai.GenerateRequest{Prompt: prompt})
	if err != nil {
		fmt.Println(err)
		return GeneratePlanResult{}, fmt.Errorf("gemini generate: %w", err)
	}

	parsed, err := ai.ParsePlanResponse(geminiResp.Text)
	if err != nil {
		return GeneratePlanResult{}, fmt.Errorf("parse plan: %w", err)
	}

	fmt.Println("this is the parsed", parsed)
	proposal, err := s.plans.CreateProposal(ctx, userID, parsed.Proposal)
	if err != nil {
		return GeneratePlanResult{}, fmt.Errorf("create proposal: %w", err)
	}
	fmt.Println("this is the proposal", proposal)

	return GeneratePlanResult{Proposal: proposal, FallbackActive: false}, nil
}

func (s *GenerateService) fallback(ctx context.Context, userID string) (GeneratePlanResult, error) {
	startDate := time.Now().UTC().Format("2006-01-02")
	fallbackPlan, err := ai.FallbackWeekPlan(startDate)
	if err != nil {
		return GeneratePlanResult{}, fmt.Errorf("fallback generation: %w", err)
	}

	proposal, err := s.plans.CreateProposal(ctx, userID, fallbackPlan.Proposal)
	if err != nil {
		return GeneratePlanResult{}, fmt.Errorf("fallback proposal: %w", err)
	}

	return GeneratePlanResult{Proposal: proposal, FallbackActive: true}, nil
}

func buildPromptContext(ctx context.Context, profile *ProfileService, pantry *PantryService, userID string, periodType string, userMessage string) ai.PromptRequest {
	req := ai.PromptRequest{
		UserRequest: userMessage,
		StartDate:   time.Now().UTC().Format("2006-01-02"),
	}

	profileResp, err := profile.GetProfile(ctx, userID)
	if err == nil {
		prefs := profileResp.Preferences
		var prefParts []string
		prefParts = append(prefParts, fmt.Sprintf("Diet type: %s", prefs.DietType))
		if len(prefs.Allergies) > 0 {
			prefParts = append(prefParts, fmt.Sprintf("Allergies: %s", strings.Join(prefs.Allergies, ", ")))
		}
		if len(prefs.Dislikes) > 0 {
			prefParts = append(prefParts, fmt.Sprintf("Dislikes: %s", strings.Join(prefs.Dislikes, ", ")))
		}
		if len(prefs.Likes) > 0 {
			prefParts = append(prefParts, fmt.Sprintf("Likes: %s", strings.Join(prefs.Likes, ", ")))
		}
		if prefs.DailyCalorieTarget != nil {
			prefParts = append(prefParts, fmt.Sprintf("Daily calorie target: %d", *prefs.DailyCalorieTarget))
		}
		req.PreferenceSummary = strings.Join(prefParts, "\n")

		metrics := profileResp.Metrics
		var metricParts []string
		if metrics.HeightCM != nil {
			metricParts = append(metricParts, fmt.Sprintf("Height: %.0f cm", *metrics.HeightCM))
		}
		if metrics.WeightKG != nil {
			metricParts = append(metricParts, fmt.Sprintf("Weight: %.0f kg", *metrics.WeightKG))
		}
		if metrics.Age != nil {
			metricParts = append(metricParts, fmt.Sprintf("Age: %d", *metrics.Age))
		}
		if metrics.Sex != nil {
			metricParts = append(metricParts, fmt.Sprintf("Sex: %s", *metrics.Sex))
		}
		if metrics.ActivityLevel != nil {
			metricParts = append(metricParts, fmt.Sprintf("Activity level: %s", *metrics.ActivityLevel))
		}
		if metrics.Goal != nil {
			metricParts = append(metricParts, fmt.Sprintf("Goal: %s", *metrics.Goal))
		}
		req.BodyMetricsSummary = strings.Join(metricParts, "\n")

		budget := profileResp.Budget
		if budget.AmountCents > 0 {
			req.BudgetTarget = fmt.Sprintf("Monthly budget: $%.2f %s", float64(budget.AmountCents)/100, budget.Currency)
		}
	}

	pantryItems, err := pantry.ListPantryItems(ctx, userID)
	if err == nil && len(pantryItems) > 0 {
		snapshot := make([]string, 0, len(pantryItems))
		for _, item := range pantryItems {
			snapshot = append(snapshot, fmt.Sprintf("%s: %.1f %s", item.Food.Description, item.Quantity, item.Unit))
		}
		req.PantrySnapshot = snapshot
	}

	return req
}

func aiPromptTemplate(periodType string) ai.PromptTemplate {
	switch periodType {
	case "meal":
		return ai.PromptTemplateSingleMeal
	case "day":
		return ai.PromptTemplateDayPlan
	case "week":
		return ai.PromptTemplateWeekPlan
	case "month":
		return ai.PromptTemplateMonthPlan
	default:
		return ai.PromptTemplateWeekPlan
	}
}

func responseContract(periodType string) string {
	return `{
  "period_type": "meal|day|week|month",
  "start_date": "YYYY-MM-DD",
  "end_date": "YYYY-MM-DD",
  "estimated_total_cost": 0.00,
  "days": [{ "date": "YYYY-MM-DD", "meals": [{ "meal_type": "breakfast|lunch|dinner|snacks", "recipe_name": "", "servings": 1, "ingredients": [{ "name": "", "quantity": 0, "unit": "" }], "macros": { "protein": 0, "carbs": 0, "fat": 0, "calories": 0 }, "estimated_cost": 0.00 }] }]
}`
}
