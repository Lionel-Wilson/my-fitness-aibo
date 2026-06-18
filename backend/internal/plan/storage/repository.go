package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/plan/domain"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

type PlanRepository interface {
	CreatePlan(ctx context.Context, userID uuid.UUID, input domain.PlanInput) (entity.Plan, error)
	ListPlans(ctx context.Context, userID uuid.UUID) ([]entity.Plan, error)
	GetPlan(ctx context.Context, userID, id uuid.UUID) (entity.Plan, error)
	UpdatePlan(ctx context.Context, userID, id uuid.UUID, input domain.PlanInput) (entity.Plan, error)
	DeletePlan(ctx context.Context, userID, id uuid.UUID) error
}

type planRepository struct {
	pool *pgxpool.Pool
}

func NewPlanRepository(pool *pgxpool.Pool) PlanRepository {
	return &planRepository{pool: pool}
}

const planCols = `id, user_id, name, quality, description, cycle_label, period_start, period_end, is_active, created_at`

func scanPlan(row pgx.Row) (entity.Plan, error) {
	var p entity.Plan
	err := row.Scan(&p.ID, &p.UserID, &p.Name, &p.Quality, &p.Description,
		&p.CycleLabel, &p.PeriodStart, &p.PeriodEnd, &p.IsActive, &p.CreatedAt)

	return p, err
}

func (r *planRepository) CreatePlan(ctx context.Context, userID uuid.UUID, input domain.PlanInput) (entity.Plan, error) {
	return scanPlan(r.pool.QueryRow(ctx,
		`INSERT INTO plans (user_id, name, quality, description, cycle_label, period_start, period_end, is_active)
		 VALUES ($1,$2,$3,$4,$5,$6,$7,$8) RETURNING `+planCols,
		userID, input.Name, input.Quality, input.Description, input.CycleLabel, input.PeriodStart, input.PeriodEnd, input.IsActive))
}

func (r *planRepository) ListPlans(ctx context.Context, userID uuid.UUID) ([]entity.Plan, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+planCols+` FROM plans WHERE user_id = $1 ORDER BY is_active DESC, created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	plans := []entity.Plan{}
	for rows.Next() {
		p, err := scanPlan(rows)
		if err != nil {
			return nil, err
		}
		plans = append(plans, p)
	}

	return plans, rows.Err()
}

func (r *planRepository) GetPlan(ctx context.Context, userID, id uuid.UUID) (entity.Plan, error) {
	p, err := scanPlan(r.pool.QueryRow(ctx,
		`SELECT `+planCols+` FROM plans WHERE id = $1 AND user_id = $2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Plan{}, commonErrors.ErrNotFound
	}

	return p, err
}

func (r *planRepository) UpdatePlan(ctx context.Context, userID, id uuid.UUID, input domain.PlanInput) (entity.Plan, error) {
	updated, err := scanPlan(r.pool.QueryRow(ctx,
		`UPDATE plans SET name=$3, quality=$4, description=$5, cycle_label=$6,
		   period_start=$7, period_end=$8, is_active=$9
		 WHERE id=$1 AND user_id=$2 RETURNING `+planCols,
		id, userID, input.Name, input.Quality, input.Description, input.CycleLabel, input.PeriodStart, input.PeriodEnd, input.IsActive))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Plan{}, commonErrors.ErrNotFound
	}

	return updated, err
}

func (r *planRepository) DeletePlan(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx, `DELETE FROM plans WHERE id=$1 AND user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return commonErrors.ErrNotFound
	}

	return nil
}
