package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"pantrypal/backend/internal/platform/id"
	"pantrypal/backend/internal/transport/http/dto"
)

type PantryRepository struct {
	db *sql.DB
}

func NewPantryRepository(db *sql.DB) *PantryRepository {
	return &PantryRepository{db: db}
}

func (r *PantryRepository) ListByUserID(ctx context.Context, userID string) ([]dto.PantryItemResponse, error) {
	rows, err := r.db.QueryContext(
		ctx,
		`SELECT p.id, p.quantity, p.unit, f.fdc_id, f.description, COALESCE(f.food_class, '')
		 FROM pantry_items p
		 JOIN usda_foods f ON f.fdc_id = p.fdc_id
		 WHERE p.user_id = ?
		 ORDER BY f.description ASC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dto.PantryItemResponse, 0)
	for rows.Next() {
		var item dto.PantryItemResponse
		if err := rows.Scan(&item.ID, &item.Quantity, &item.Unit, &item.Food.FDCID, &item.Food.Description, &item.Food.FoodClass); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *PantryRepository) Upsert(ctx context.Context, userID string, req dto.PantryItemRequest) (dto.PantryItemResponse, error) {
	itemID, err := id.New("pnt")
	if err != nil {
		return dto.PantryItemResponse{}, err
	}

	unit := strings.TrimSpace(req.Unit)
	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO pantry_items (id, user_id, fdc_id, quantity, unit)
		 VALUES (?, ?, ?, ?, ?)
		 ON CONFLICT(user_id, fdc_id, unit) DO UPDATE SET
			 quantity = pantry_items.quantity + excluded.quantity,
			 updated_at = CURRENT_TIMESTAMP`,
		itemID,
		userID,
		req.FDCID,
		req.Quantity,
		unit,
	)
	if err != nil {
		return dto.PantryItemResponse{}, err
	}

	return r.GetByUserAndFood(ctx, userID, req.FDCID, unit)
}

func (r *PantryRepository) GetByUserAndFood(ctx context.Context, userID string, fdcID int64, unit string) (dto.PantryItemResponse, error) {
	var item dto.PantryItemResponse
	err := r.db.QueryRowContext(
		ctx,
		`SELECT p.id, p.quantity, p.unit, f.fdc_id, f.description, COALESCE(f.food_class, '')
		 FROM pantry_items p
		 JOIN usda_foods f ON f.fdc_id = p.fdc_id
		 WHERE p.user_id = ? AND p.fdc_id = ? AND p.unit = ?`,
		userID,
		fdcID,
		unit,
	).Scan(&item.ID, &item.Quantity, &item.Unit, &item.Food.FDCID, &item.Food.Description, &item.Food.FoodClass)
	return item, err
}

func (r *PantryRepository) PatchQuantity(ctx context.Context, userID, itemID string, delta float64) (dto.PantryItemResponse, error) {
	_, err := r.db.ExecContext(
		ctx,
		`UPDATE pantry_items
		 SET quantity = CASE
			 WHEN quantity + ? < 0 THEN 0
			 ELSE quantity + ?
		 END,
		 updated_at = CURRENT_TIMESTAMP
		 WHERE id = ? AND user_id = ?`,
		delta,
		delta,
		itemID,
		userID,
	)
	if err != nil {
		return dto.PantryItemResponse{}, err
	}

	return r.GetByID(ctx, userID, itemID)
}

func (r *PantryRepository) GetByID(ctx context.Context, userID, itemID string) (dto.PantryItemResponse, error) {
	var item dto.PantryItemResponse
	err := r.db.QueryRowContext(
		ctx,
		`SELECT p.id, p.quantity, p.unit, f.fdc_id, f.description, COALESCE(f.food_class, '')
		 FROM pantry_items p
		 JOIN usda_foods f ON f.fdc_id = p.fdc_id
		 WHERE p.id = ? AND p.user_id = ?`,
		itemID,
		userID,
	).Scan(&item.ID, &item.Quantity, &item.Unit, &item.Food.FDCID, &item.Food.Description, &item.Food.FoodClass)
	return item, err
}

func (r *PantryRepository) Delete(ctx context.Context, userID, itemID string) error {
	result, err := r.db.ExecContext(ctx, `DELETE FROM pantry_items WHERE id = ? AND user_id = ?`, itemID, userID)
	if err != nil {
		return err
	}
	affected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type PantryDeduction struct {
	PantryItemID string
	BeforeQty    float64
	AfterQty     float64
	FoodDesc     string
}

func (r *PantryRepository) DeductByFood(ctx context.Context, userID string, fdcID int64, unit string, quantity float64) (PantryDeduction, error) {
	var d PantryDeduction
	err := r.db.QueryRowContext(
		ctx,
		`SELECT p.id, p.quantity, f.description
		 FROM pantry_items p
		 JOIN usda_foods f ON f.fdc_id = p.fdc_id
		 WHERE p.user_id = ? AND p.fdc_id = ? AND p.unit = ?`,
		userID,
		fdcID,
		unit,
	).Scan(&d.PantryItemID, &d.BeforeQty, &d.FoodDesc)
	if err != nil {
		return PantryDeduction{}, err
	}

	newQty := d.BeforeQty - quantity
	if newQty < 0 {
		newQty = 0
	}

	_, err = r.db.ExecContext(
		ctx,
		`UPDATE pantry_items SET quantity = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		newQty,
		d.PantryItemID,
	)
	if err != nil {
		return PantryDeduction{}, err
	}

	d.AfterQty = newQty
	return d, nil
}

func (r *PantryRepository) Exists(ctx context.Context, userID, itemID string) (bool, error) {
	_, err := r.GetByID(ctx, userID, itemID)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
