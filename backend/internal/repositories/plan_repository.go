package repositories

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"pantrypal/backend/internal/platform/id"
	"pantrypal/backend/internal/transport/http/dto"
)

type StoredPlan struct {
	ID               string
	UserID           string
	PeriodType       string
	StartDate        string
	EndDate          string
	Status           string
	Source           string
	ProposalVersion  int
	AICostCentsTotal int
	Notes            string
	Meals            []StoredPlanMeal
}

type StoredPlanMeal struct {
	ID                 string
	RecipeID           string
	RecipeName         string
	ScheduledDate      string
	MealSection        string
	Servings           float64
	EstimatedCostCents int
	Calories           float64
	ProteinG           float64
	CarbsG             float64
	FatG               float64
	IsConsumed         bool
	ConsumedAt         string
}

type PlanRepository struct {
	db *sql.DB
}

func NewPlanRepository(db *sql.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) CreateProposal(ctx context.Context, userID string, req dto.PlanProposalRequest) (StoredPlan, error) {
	planID, err := id.New("pln")
	if err != nil {
		return StoredPlan{}, err
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return StoredPlan{}, err
	}
	defer tx.Rollback()

	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO meal_plans (id, user_id, period_type, start_date, end_date, status, source, proposal_version, ai_cost_cents_total, notes)
		 VALUES (?, ?, ?, ?, ?, 'proposal', ?, ?, ?, ?)`,
		planID,
		userID,
		req.PeriodType,
		req.StartDate,
		req.EndDate,
		req.Source,
		req.ProposalVersion,
		req.AICostCentsTotal,
		req.Notes,
	)
	if err != nil {
		return StoredPlan{}, err
	}

	for _, meal := range req.Meals {
		mealID, err := id.New("pm")
		if err != nil {
			return StoredPlan{}, err
		}

		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO plan_meals (
				id, meal_plan_id, recipe_id, scheduled_date, meal_section, recipe_name,
				servings, kcal, protein_g, carbs_g, fat_g, estimated_cost_cents
			 ) VALUES (?, ?, NULLIF(?, ''), ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			mealID,
			planID,
			meal.RecipeID,
			meal.ScheduledDate,
			meal.MealSection,
			meal.RecipeName,
			meal.Servings,
			meal.Macros.Calories,
			meal.Macros.ProteinG,
			meal.Macros.CarbsG,
			meal.Macros.FatG,
			meal.EstimatedCostCents,
		)
		if err != nil {
			return StoredPlan{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return StoredPlan{}, err
	}

	return r.GetPlanByID(ctx, userID, planID)
}

func (r *PlanRepository) GetPlanByID(ctx context.Context, userID, planID string) (StoredPlan, error) {
	var plan StoredPlan
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, user_id, period_type, start_date, end_date, status, source, proposal_version, ai_cost_cents_total, COALESCE(notes, '')
		 FROM meal_plans
		 WHERE id = ? AND user_id = ?`,
		planID,
		userID,
	).Scan(
		&plan.ID,
		&plan.UserID,
		&plan.PeriodType,
		&plan.StartDate,
		&plan.EndDate,
		&plan.Status,
		&plan.Source,
		&plan.ProposalVersion,
		&plan.AICostCentsTotal,
		&plan.Notes,
	)
	if err != nil {
		return StoredPlan{}, err
	}
	plan.StartDate = normalizeDateOnly(plan.StartDate)
	plan.EndDate = normalizeDateOnly(plan.EndDate)

	meals, err := r.listMealsForPlan(ctx, planID)
	if err != nil {
		return StoredPlan{}, err
	}
	plan.Meals = meals
	return plan, nil
}

func (r *PlanRepository) UpdateStatus(ctx context.Context, userID, planID, status, noteSuffix string) (StoredPlan, error) {
	if strings.TrimSpace(noteSuffix) == "" {
		_, err := r.db.ExecContext(
			ctx,
			`UPDATE meal_plans SET status = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ? AND user_id = ?`,
			status,
			planID,
			userID,
		)
		if err != nil {
			return StoredPlan{}, err
		}
	} else {
		_, err := r.db.ExecContext(
			ctx,
			`UPDATE meal_plans
			 SET status = ?,
			     notes = TRIM(COALESCE(notes, '') || CASE WHEN COALESCE(notes, '') = '' THEN '' ELSE ' | ' END || ?),
			     updated_at = CURRENT_TIMESTAMP
			 WHERE id = ? AND user_id = ?`,
			status,
			noteSuffix,
			planID,
			userID,
		)
		if err != nil {
			return StoredPlan{}, err
		}
	}

	return r.GetPlanByID(ctx, userID, planID)
}

func (r *PlanRepository) ListAcceptedPlansForRange(ctx context.Context, userID, startDate, endDate string) ([]StoredPlan, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT id, user_id, period_type, start_date, end_date, status, source, proposal_version, ai_cost_cents_total, COALESCE(notes, '')
		 FROM meal_plans
		 WHERE user_id = ? AND status = 'accepted' AND start_date <= ? AND end_date >= ?
		 ORDER BY start_date ASC, created_at ASC`,
		userID,
		endDate,
		startDate,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := make([]StoredPlan, 0)
	for rows.Next() {
		var plan StoredPlan
		if err := rows.Scan(
			&plan.ID,
			&plan.UserID,
			&plan.PeriodType,
			&plan.StartDate,
			&plan.EndDate,
			&plan.Status,
			&plan.Source,
			&plan.ProposalVersion,
			&plan.AICostCentsTotal,
			&plan.Notes,
		); err != nil {
			return nil, err
		}
		plan.StartDate = normalizeDateOnly(plan.StartDate)
		plan.EndDate = normalizeDateOnly(plan.EndDate)
		plans = append(plans, plan)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	for i := range plans {
		meals, err := r.listMealsForPlanInRange(ctx, plans[i].ID, startDate, endDate)
		if err != nil {
			return nil, err
		}
		plans[i].Meals = meals
	}

	return plans, nil
}

func (r *PlanRepository) GetPlanMealByID(ctx context.Context, userID, mealID string) (StoredPlanMeal, error) {
	var meal StoredPlanMeal
	var consumedInt int
	err := r.db.QueryRowContext(
		ctx,
		`SELECT pm.id, COALESCE(pm.recipe_id, ''), pm.recipe_name, pm.scheduled_date, pm.meal_section,
		        pm.servings, COALESCE(pm.kcal, 0), COALESCE(pm.protein_g, 0), COALESCE(pm.carbs_g, 0),
		        COALESCE(pm.fat_g, 0), pm.estimated_cost_cents, pm.is_consumed, COALESCE(pm.consumed_at, '')
		 FROM plan_meals pm
		 JOIN meal_plans mp ON mp.id = pm.meal_plan_id
		 WHERE pm.id = ? AND mp.user_id = ?`,
		mealID,
		userID,
	).Scan(
		&meal.ID,
		&meal.RecipeID,
		&meal.RecipeName,
		&meal.ScheduledDate,
		&meal.MealSection,
		&meal.Servings,
		&meal.Calories,
		&meal.ProteinG,
		&meal.CarbsG,
		&meal.FatG,
		&meal.EstimatedCostCents,
		&consumedInt,
		&meal.ConsumedAt,
	)
	if err != nil {
		return StoredPlanMeal{}, err
	}
	meal.ScheduledDate = normalizeDateOnly(meal.ScheduledDate)
	meal.IsConsumed = consumedInt == 1
	return meal, nil
}

func (r *PlanRepository) MarkMealConsumed(ctx context.Context, mealID string) error {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE plan_meals SET is_consumed = 1, consumed_at = CURRENT_TIMESTAMP, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		mealID,
	)
	return err
}

func (r *PlanRepository) listMealsForPlan(ctx context.Context, planID string) ([]StoredPlanMeal, error) {
	return r.listMeals(ctx, `WHERE meal_plan_id = ?`, planID)
}

func (r *PlanRepository) listMealsForPlanInRange(ctx context.Context, planID, startDate, endDate string) ([]StoredPlanMeal, error) {
	return r.listMeals(ctx, `WHERE meal_plan_id = ? AND scheduled_date >= ? AND scheduled_date <= ?`, planID, startDate, endDate)
}

func (r *PlanRepository) listMeals(ctx context.Context, whereClause string, args ...any) ([]StoredPlanMeal, error) {
	query := fmt.Sprintf(`SELECT id, COALESCE(recipe_id, ''), recipe_name, scheduled_date, meal_section, servings,
		 COALESCE(kcal, 0), COALESCE(protein_g, 0), COALESCE(carbs_g, 0), COALESCE(fat_g, 0), estimated_cost_cents,
		 is_consumed, COALESCE(consumed_at, '')
		 FROM plan_meals %s ORDER BY scheduled_date ASC, meal_section ASC`, whereClause)

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	meals := make([]StoredPlanMeal, 0)
	for rows.Next() {
		var meal StoredPlanMeal
		var consumedInt int
		if err := rows.Scan(
			&meal.ID,
			&meal.RecipeID,
			&meal.RecipeName,
			&meal.ScheduledDate,
			&meal.MealSection,
			&meal.Servings,
			&meal.Calories,
			&meal.ProteinG,
			&meal.CarbsG,
			&meal.FatG,
			&meal.EstimatedCostCents,
			&consumedInt,
			&meal.ConsumedAt,
		); err != nil {
			return nil, err
		}
		meal.ScheduledDate = normalizeDateOnly(meal.ScheduledDate)
		meal.IsConsumed = consumedInt == 1
		meals = append(meals, meal)
	}

	return meals, rows.Err()
}

func normalizeDateOnly(value string) string {
	value = strings.TrimSpace(value)
	if value == "" {
		return value
	}
	if len(value) >= 10 && value[4] == '-' && value[7] == '-' {
		return value[:10]
	}
	for _, layout := range []string{time.RFC3339, "2006-01-02 15:04:05Z07:00", "2006-01-02 15:04:05", "2006-01-02"} {
		if t, err := time.Parse(layout, value); err == nil {
			return t.Format("2006-01-02")
		}
	}
	return value
}
