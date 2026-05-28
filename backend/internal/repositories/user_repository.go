package repositories

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"pantrypal/backend/internal/platform/id"
	"pantrypal/backend/internal/transport/http/dto"
)

type User struct {
	ID           string
	Email        string
	PasswordHash string
	DisplayName  string
}

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) ExistsByEmail(ctx context.Context, email string) (bool, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM users WHERE email = ?`, email).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *UserRepository) CreateUser(ctx context.Context, email, passwordHash, displayName string) (User, error) {
	userID, err := id.New("usr")
	if err != nil {
		return User{}, err
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO users (id, email, password_hash, display_name) VALUES (?, ?, ?, ?)`,
		userID,
		email,
		passwordHash,
		displayName,
	)
	if err != nil {
		return User{}, err
	}

	if err := r.EnsureDefaultProfileRows(ctx, userID); err != nil {
		return User{}, err
	}

	return User{ID: userID, Email: email, DisplayName: displayName}, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (User, error) {
	var out User
	var displayName sql.NullString
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, password_hash, display_name FROM users WHERE email = ?`,
		email,
	).Scan(&out.ID, &out.Email, &out.PasswordHash, &displayName)
	if err != nil {
		return User{}, err
	}
	if displayName.Valid {
		out.DisplayName = displayName.String
	}
	return out, nil
}

func (r *UserRepository) GetByID(ctx context.Context, userID string) (User, error) {
	var out User
	err := r.db.QueryRowContext(
		ctx,
		`SELECT id, email, COALESCE(display_name, '') FROM users WHERE id = ?`,
		userID,
	).Scan(&out.ID, &out.Email, &out.DisplayName)
	if err != nil {
		return User{}, err
	}
	return out, nil
}

func (r *UserRepository) ExistsByID(ctx context.Context, userID string) (bool, error) {
	var n int
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(1) FROM users WHERE id = ?`, userID).Scan(&n); err != nil {
		return false, err
	}
	return n > 0, nil
}

func (r *UserRepository) EnsureDefaultProfileRows(ctx context.Context, userID string) error {
	if _, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user_preferences (id, user_id, diet_type, allergies_json, dislikes_json, likes_json)
         VALUES (?, ?, 'omnivore', '[]', '[]', '[]')
         ON CONFLICT(user_id) DO NOTHING`,
		id.Must("upr"),
		userID,
	); err != nil {
		return err
	}
	if _, err := r.db.ExecContext(
		ctx,
		`INSERT INTO user_body_metrics (id, user_id)
         VALUES (?, ?)
         ON CONFLICT(user_id) DO NOTHING`,
		id.Must("ubm"),
		userID,
	); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetProfile(ctx context.Context, userID string) (dto.ProfileResponse, error) {
	if err := r.EnsureDefaultProfileRows(ctx, userID); err != nil {
		return dto.ProfileResponse{}, err
	}

	user, err := r.GetByID(ctx, userID)
	if err != nil {
		return dto.ProfileResponse{}, err
	}

	out := dto.ProfileResponse{
		User: dto.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			DisplayName: user.DisplayName,
		},
	}

	var sex sql.NullString
	var activity sql.NullString
	var goal sql.NullString
	err = r.db.QueryRowContext(
		ctx,
		`SELECT height_cm, weight_kg, age, sex, activity_level, goal FROM user_body_metrics WHERE user_id = ?`,
		userID,
	).Scan(&out.Metrics.HeightCM, &out.Metrics.WeightKG, &out.Metrics.Age, &sex, &activity, &goal)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return dto.ProfileResponse{}, err
	}
	if sex.Valid {
		out.Metrics.Sex = &sex.String
	}
	if activity.Valid {
		out.Metrics.ActivityLevel = &activity.String
	}
	if goal.Valid {
		out.Metrics.Goal = &goal.String
	}

	var allergiesJSON string
	var dislikesJSON string
	var likesJSON string
	var notes sql.NullString
	err = r.db.QueryRowContext(
		ctx,
		`SELECT diet_type, allergies_json, dislikes_json, likes_json, daily_calorie_target, notes
         FROM user_preferences WHERE user_id = ?`,
		userID,
	).Scan(&out.Preferences.DietType, &allergiesJSON, &dislikesJSON, &likesJSON, &out.Preferences.DailyCalorieTarget, &notes)
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return dto.ProfileResponse{}, err
	}
	if out.Preferences.DietType == "" {
		out.Preferences.DietType = "omnivore"
	}
	_ = json.Unmarshal([]byte(allergiesJSON), &out.Preferences.Allergies)
	_ = json.Unmarshal([]byte(dislikesJSON), &out.Preferences.Dislikes)
	_ = json.Unmarshal([]byte(likesJSON), &out.Preferences.Likes)
	if notes.Valid {
		out.Preferences.Notes = &notes.String
	}

	err = r.db.QueryRowContext(
		ctx,
		`SELECT month, currency, amount_cents FROM budgets WHERE user_id = ? ORDER BY month DESC LIMIT 1`,
		userID,
	).Scan(&out.Budget.Month, &out.Budget.Currency, &out.Budget.AmountCents)
	if errors.Is(err, sql.ErrNoRows) {
		out.Budget.Month = time.Now().UTC().Format("2006-01")
		out.Budget.Currency = "USD"
		out.Budget.AmountCents = 0
		return out, nil
	}
	if err != nil {
		return dto.ProfileResponse{}, err
	}
	return out, nil
}

func (r *UserRepository) UpsertMetrics(ctx context.Context, userID string, req dto.PatchMetricsRequest) error {
	rowID, err := id.New("ubm")
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO user_body_metrics (id, user_id, height_cm, weight_kg, age, sex, activity_level, goal)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)
         ON CONFLICT(user_id) DO UPDATE SET
             height_cm = excluded.height_cm,
             weight_kg = excluded.weight_kg,
             age = excluded.age,
             sex = excluded.sex,
             activity_level = excluded.activity_level,
             goal = excluded.goal,
             updated_at = CURRENT_TIMESTAMP`,
		rowID,
		userID,
		req.HeightCM,
		req.WeightKG,
		req.Age,
		req.Sex,
		req.ActivityLevel,
		req.Goal,
	)
	return err
}

func (r *UserRepository) UpsertPreferences(ctx context.Context, userID string, req dto.PatchPreferencesRequest) error {
	rowID, err := id.New("upr")
	if err != nil {
		return err
	}
	allergiesJSON, _ := json.Marshal(req.Allergies)
	dislikesJSON, _ := json.Marshal(req.Dislikes)
	likesJSON, _ := json.Marshal(req.Likes)

	dietType := "omnivore"
	if req.DietType != nil && strings.TrimSpace(*req.DietType) != "" {
		dietType = strings.TrimSpace(*req.DietType)
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO user_preferences (id, user_id, diet_type, allergies_json, dislikes_json, likes_json, daily_calorie_target, notes)
         VALUES (?, ?, ?, ?, ?, ?, ?, ?)
         ON CONFLICT(user_id) DO UPDATE SET
             diet_type = excluded.diet_type,
             allergies_json = excluded.allergies_json,
             dislikes_json = excluded.dislikes_json,
             likes_json = excluded.likes_json,
             daily_calorie_target = excluded.daily_calorie_target,
             notes = excluded.notes,
             updated_at = CURRENT_TIMESTAMP`,
		rowID,
		userID,
		dietType,
		string(allergiesJSON),
		string(dislikesJSON),
		string(likesJSON),
		req.DailyCalorieTarget,
		req.Notes,
	)
	return err
}

func (r *UserRepository) UpsertBudget(ctx context.Context, userID string, req dto.PatchBudgetRequest) error {
	rowID, err := id.New("bdg")
	if err != nil {
		return err
	}

	month := time.Now().UTC().Format("2006-01")
	if req.Month != nil && strings.TrimSpace(*req.Month) != "" {
		month = strings.TrimSpace(*req.Month)
	}
	currency := "USD"
	if req.Currency != nil && strings.TrimSpace(*req.Currency) != "" {
		currency = strings.ToUpper(strings.TrimSpace(*req.Currency))
	}

	_, err = r.db.ExecContext(
		ctx,
		`INSERT INTO budgets (id, user_id, month, currency, amount_cents)
         VALUES (?, ?, ?, ?, ?)
         ON CONFLICT(user_id, month) DO UPDATE SET
             currency = excluded.currency,
             amount_cents = excluded.amount_cents,
             updated_at = CURRENT_TIMESTAMP`,
		rowID,
		userID,
		month,
		currency,
		*req.AmountCents,
	)
	return err
}
