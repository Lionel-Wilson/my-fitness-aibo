package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Plan is a training block, named after the quality it trains.
type Plan struct {
	ID          uuid.UUID  `json:"id"`
	UserID      uuid.UUID  `json:"-"`
	Name        string     `json:"name"`
	Quality     string     `json:"quality"`
	Description string     `json:"description"`
	CycleLabel  string     `json:"cycleLabel"`
	PeriodStart *time.Time `json:"periodStart"`
	PeriodEnd   *time.Time `json:"periodEnd"`
	IsActive    bool       `json:"isActive"`
	CreatedAt   time.Time  `json:"createdAt"`
}

// PlanParams holds the editable fields of a plan.
type PlanParams struct {
	Name        string
	Quality     string
	Description string
	CycleLabel  string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
	IsActive    bool
}

const planCols = `id, user_id, name, quality, description, cycle_label, period_start, period_end, is_active, created_at`

func scanPlan(row pgx.Row) (Plan, error) {
	var p Plan
	err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.Quality, &p.Description,
		&p.CycleLabel, &p.PeriodStart, &p.PeriodEnd, &p.IsActive, &p.CreatedAt)
	return p, err
}

// CreatePlan inserts a plan for the user.
func (s *Store) CreatePlan(ctx context.Context, userID uuid.UUID, p PlanParams) (Plan, error) {
	return scanPlan(s.pool.QueryRow(ctx,
		`INSERT INTO plans (user_id, name, quality, description, cycle_label, period_start, period_end, is_active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+planCols,
		userID, p.Name, p.Quality, p.Description, p.CycleLabel, p.PeriodStart, p.PeriodEnd, p.IsActive))
}

// ListPlans returns the user's plans, newest first.
func (s *Store) ListPlans(ctx context.Context, userID uuid.UUID) ([]Plan, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+planCols+` FROM plans WHERE user_id = $1 ORDER BY is_active DESC, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []Plan{}
	for rows.Next() {
		p, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}
	return plans, rows.Err()
}

// GetPlan returns a single plan owned by the user.
func (s *Store) GetPlan(ctx context.Context, userID, id uuid.UUID) (Plan, error) {
	p, err := scanPlan(s.pool.QueryRow(ctx,
		`SELECT `+planCols+` FROM plans WHERE id = $1 AND user_id = $2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Plan{}, ErrNotFound
	}
	return p, err
}

// UpdatePlan updates a plan owned by the user.
func (s *Store) UpdatePlan(ctx context.Context, userID, id uuid.UUID, p PlanParams) (Plan, error) {
	updated, err := scanPlan(s.pool.QueryRow(ctx,
		`UPDATE plans SET name=$3, quality=$4, description=$5, cycle_label=$6,
		   period_start=$7, period_end=$8, is_active=$9
		 WHERE id=$1 AND user_id=$2 RETURNING `+planCols,
		id, userID, p.Name, p.Quality, p.Description, p.CycleLabel, p.PeriodStart, p.PeriodEnd, p.IsActive))
	if errors.Is(err, pgx.ErrNoRows) {
		return Plan{}, ErrNotFound
	}
	return updated, err
}

// DeletePlan removes a plan owned by the user (cascades to children).
func (s *Store) DeletePlan(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx, `DELETE FROM plans WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
