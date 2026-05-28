package services

import (
	"context"
	"database/sql"
	"errors"
	"sort"
	"strings"
	"time"

	"pantrypal/backend/internal/repositories"
	"pantrypal/backend/internal/transport/http/dto"
)

var (
	ErrInvalidPlanType    = errors.New("periodType must be one of day, week, month")
	ErrInvalidStartDate   = errors.New("startDate must use YYYY-MM-DD")
	ErrInvalidEndDate     = errors.New("endDate must use YYYY-MM-DD")
	ErrInvalidPlanMeals   = errors.New("at least one meal is required")
	ErrInvalidMealSection = errors.New("mealSection must be one of breakfast, lunch, dinner, snacks")
	ErrInvalidPlanStatus  = errors.New("plan not found")
	ErrPlanNotProposal    = errors.New("plan is not in proposal status")
	ErrWeekStartDate      = errors.New("start must use YYYY-MM-DD")
)

type PlanService struct {
	plans *repositories.PlanRepository
}

func NewPlanService(plans *repositories.PlanRepository) *PlanService {
	return &PlanService{plans: plans}
}

func (s *PlanService) CreateProposal(ctx context.Context, userID string, req dto.PlanProposalRequest) (dto.ProposalResponse, error) {
	if err := validatePlanProposal(req); err != nil {
		return dto.ProposalResponse{}, err
	}
	if req.Source == "" {
		req.Source = "ai"
	}
	if req.ProposalVersion == 0 {
		req.ProposalVersion = 1
	}
	if req.EndDate == "" {
		req.EndDate = req.StartDate
	}

	plan, err := s.plans.CreateProposal(ctx, userID, req)
	if err != nil {
		return dto.ProposalResponse{}, err
	}

	return buildProposalResponse([]repositories.StoredPlan{plan}), nil
}

func (s *PlanService) AcceptProposal(ctx context.Context, userID, planID string) (dto.ProposalResponse, error) {
	plan, err := s.plans.GetPlanByID(ctx, userID, planID)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.ProposalResponse{}, ErrInvalidPlanStatus
	}
	if err != nil {
		return dto.ProposalResponse{}, err
	}
	if plan.Status != "proposal" {
		return dto.ProposalResponse{}, ErrPlanNotProposal
	}

	updated, err := s.plans.UpdateStatus(ctx, userID, planID, "accepted", "")
	if err != nil {
		return dto.ProposalResponse{}, err
	}

	return buildProposalResponse([]repositories.StoredPlan{updated}), nil
}

func (s *PlanService) DeclineProposal(ctx context.Context, userID, planID, reason string) (dto.ProposalResponse, error) {
	plan, err := s.plans.GetPlanByID(ctx, userID, planID)
	if errors.Is(err, sql.ErrNoRows) {
		return dto.ProposalResponse{}, ErrInvalidPlanStatus
	}
	if err != nil {
		return dto.ProposalResponse{}, err
	}
	if plan.Status != "proposal" {
		return dto.ProposalResponse{}, ErrPlanNotProposal
	}

	note := ""
	if strings.TrimSpace(reason) != "" {
		note = "decline_reason: " + strings.TrimSpace(reason)
	}

	updated, err := s.plans.UpdateStatus(ctx, userID, planID, "declined", note)
	if err != nil {
		return dto.ProposalResponse{}, err
	}

	return buildProposalResponse([]repositories.StoredPlan{updated}), nil
}

func (s *PlanService) GetWeekPlan(ctx context.Context, userID, startDate string) (dto.WeekPlanResponse, error) {
	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		return dto.WeekPlanResponse{}, ErrWeekStartDate
	}
	endDate := start.AddDate(0, 0, 6).Format("2006-01-02")
	plans, err := s.plans.ListAcceptedPlansForRange(ctx, userID, startDate, endDate)
	if err != nil {
		return dto.WeekPlanResponse{}, err
	}
	return buildWeekPlanResponse(plans, startDate, endDate), nil
}

func validatePlanProposal(req dto.PlanProposalRequest) error {
	switch req.PeriodType {
	case "day", "week", "month":
	default:
		return ErrInvalidPlanType
	}
	if _, err := time.Parse("2006-01-02", req.StartDate); err != nil {
		return ErrInvalidStartDate
	}
	if req.EndDate != "" {
		if _, err := time.Parse("2006-01-02", req.EndDate); err != nil {
			return ErrInvalidEndDate
		}
	}
	if len(req.Meals) == 0 {
		return ErrInvalidPlanMeals
	}
	for _, meal := range req.Meals {
		switch meal.MealSection {
		case "breakfast", "lunch", "dinner", "snacks":
		default:
			return ErrInvalidMealSection
		}
	}
	return nil
}

func buildProposalResponse(plans []repositories.StoredPlan) dto.ProposalResponse {
	if len(plans) == 0 {
		return dto.ProposalResponse{}
	}
	week := buildWeekPlanResponse(plans, plans[0].StartDate, plans[0].EndDate)
	if len(week.Plans) == 0 {
		return dto.ProposalResponse{}
	}
	return dto.ProposalResponse{
		Plan:       week.Plans[0],
		Days:       week.Days,
		WeekTotals: week.WeekTotals,
	}
}

func buildWeekPlanResponse(plans []repositories.StoredPlan, startDate, endDate string) dto.WeekPlanResponse {
	planSummaries := make([]dto.PlanSummaryResponse, 0, len(plans))
	allMeals := make([]repositories.StoredPlanMeal, 0)
	for _, plan := range plans {
		planSummaries = append(planSummaries, dto.PlanSummaryResponse{
			ID:               plan.ID,
			PeriodType:       plan.PeriodType,
			StartDate:        plan.StartDate,
			EndDate:          plan.EndDate,
			Status:           plan.Status,
			Source:           plan.Source,
			ProposalVersion:  plan.ProposalVersion,
			AICostCentsTotal: plan.AICostCentsTotal,
			Notes:            plan.Notes,
		})
		allMeals = append(allMeals, plan.Meals...)
	}

	sort.Slice(allMeals, func(i, j int) bool {
		if allMeals[i].ScheduledDate == allMeals[j].ScheduledDate {
			return allMeals[i].MealSection < allMeals[j].MealSection
		}
		return allMeals[i].ScheduledDate < allMeals[j].ScheduledDate
	})

	days := make([]dto.PlanDayResponse, 0)
	weeklyTotals := dto.RecipeMacrosResponse{}
	mealIndex := map[string]repositories.StoredPlanMeal{}
	for _, meal := range allMeals {
		mealIndex[meal.ScheduledDate+"|"+meal.MealSection] = meal
		weeklyTotals.Calories += meal.Calories
		weeklyTotals.ProteinG += meal.ProteinG
		weeklyTotals.CarbsG += meal.CarbsG
		weeklyTotals.FatG += meal.FatG
	}

	start, err := time.Parse("2006-01-02", startDate)
	if err != nil {
		start = deriveStartDate(plans, startDate)
	}
	end, err := time.Parse("2006-01-02", endDate)
	if err != nil {
		end = deriveEndDate(plans, endDate, start)
	}

	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		date := d.Format("2006-01-02")
		day := dto.PlanDayResponse{Date: date}
		for _, section := range []string{"breakfast", "lunch", "dinner", "snacks"} {
			meal, ok := mealIndex[date+"|"+section]
			if !ok {
				continue
			}
			res := mealToDTO(meal)
			switch section {
			case "breakfast":
				day.Sections.Breakfast = &res
			case "lunch":
				day.Sections.Lunch = &res
			case "dinner":
				day.Sections.Dinner = &res
			case "snacks":
				day.Sections.Snacks = &res
			}
			day.Totals.Calories += meal.Calories
			day.Totals.ProteinG += meal.ProteinG
			day.Totals.CarbsG += meal.CarbsG
			day.Totals.FatG += meal.FatG
		}
		days = append(days, day)
	}

	return dto.WeekPlanResponse{Plans: planSummaries, Days: days, WeekTotals: weeklyTotals}
}

func deriveStartDate(plans []repositories.StoredPlan, fallback string) time.Time {
	if fallback != "" {
		if t, err := time.Parse("2006-01-02", fallback); err == nil {
			return t
		}
	}
	for _, plan := range plans {
		if t, err := time.Parse("2006-01-02", plan.StartDate); err == nil {
			return t
		}
	}
	return time.Now().UTC()
}

func deriveEndDate(plans []repositories.StoredPlan, fallback string, start time.Time) time.Time {
	if fallback != "" {
		if t, err := time.Parse("2006-01-02", fallback); err == nil {
			return t
		}
	}
	for _, plan := range plans {
		if t, err := time.Parse("2006-01-02", plan.EndDate); err == nil {
			return t
		}
	}
	return start
}

func mealToDTO(meal repositories.StoredPlanMeal) dto.PlanMealResponse {
	return dto.PlanMealResponse{
		ID:                 meal.ID,
		RecipeID:           meal.RecipeID,
		RecipeName:         meal.RecipeName,
		ScheduledDate:      meal.ScheduledDate,
		MealSection:        meal.MealSection,
		Servings:           meal.Servings,
		EstimatedCostCents: meal.EstimatedCostCents,
		Macros: dto.RecipeMacrosResponse{
			Calories: meal.Calories,
			ProteinG: meal.ProteinG,
			CarbsG:   meal.CarbsG,
			FatG:     meal.FatG,
		},
		IsConsumed: meal.IsConsumed,
		ConsumedAt: meal.ConsumedAt,
	}
}
