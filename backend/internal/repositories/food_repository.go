package repositories

import (
	"context"
	"database/sql"
	"strings"

	"pantrypal/backend/internal/transport/http/dto"
)

type FoodRepository struct {
	db *sql.DB
}

func NewFoodRepository(db *sql.DB) *FoodRepository {
	return &FoodRepository{db: db}
}

func (r *FoodRepository) SearchFoods(ctx context.Context, query string, limit int) ([]dto.FoodSearchItem, error) {
	query = strings.TrimSpace(query)
	if limit <= 0 {
		limit = 20
	}

	rows, err := r.db.QueryContext(
		ctx,
		`SELECT fdc_id, description, COALESCE(food_class, '')
		 FROM usda_foods
		 WHERE description LIKE ?
		 ORDER BY description ASC
		 LIMIT ?`,
		"%"+query+"%",
		limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]dto.FoodSearchItem, 0)
	for rows.Next() {
		var item dto.FoodSearchItem
		if err := rows.Scan(&item.FDCID, &item.Description, &item.FoodClass); err != nil {
			return nil, err
		}
		items = append(items, item)
	}

	return items, rows.Err()
}

func (r *FoodRepository) ExistsByFDCID(ctx context.Context, fdcID int64) (bool, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM usda_foods WHERE fdc_id = ?`, fdcID).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}
