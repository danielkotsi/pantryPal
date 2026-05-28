package repositories

import (
	"context"
	"database/sql"

	"pantrypal/backend/internal/platform/id"
)

type ConsumptionLogRepository struct {
	db *sql.DB
}

func NewConsumptionLogRepository(db *sql.DB) *ConsumptionLogRepository {
	return &ConsumptionLogRepository{db: db}
}

type ConsumptionLogEntry struct {
	PlanMealID       string
	UserID           string
	FDCID            int64
	PantryItemID     string
	QuantityDeducted float64
	Unit             string
	BeforeQuantity   float64
	AfterQuantity    float64
	Warning          string
}

func (r *ConsumptionLogRepository) Insert(ctx context.Context, entry ConsumptionLogEntry) error {
	logID, err := id.New("cl")
	if err != nil {
		return err
	}

	var warning *string
	if entry.Warning != "" {
		warning = &entry.Warning
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO consumption_log (id, plan_meal_id, user_id, fdc_id, pantry_item_id, quantity_deducted, unit, before_quantity, after_quantity, warning)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		logID,
		nullableStr(entry.PlanMealID),
		entry.UserID,
		entry.FDCID,
		nullableStr(entry.PantryItemID),
		entry.QuantityDeducted,
		entry.Unit,
		entry.BeforeQuantity,
		entry.AfterQuantity,
		warning,
	)
	return err
}

func nullableStr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}
