package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/cycle/domain"
	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

type CycleRepository interface {
	CreateCycle(ctx context.Context, userID, planID uuid.UUID, input domain.CreateCycleInput) (entity.Cycle, error)
	ListCycles(ctx context.Context, userID, planID uuid.UUID) ([]entity.Cycle, error)
	UpdateCycle(ctx context.Context, userID, id uuid.UUID, input domain.UpdateCycleInput) (entity.Cycle, error)
	DeleteCycle(ctx context.Context, userID, id uuid.UUID) error
}

type cycleRepository struct {
	pool *pgxpool.Pool
}

func NewCycleRepository(pool *pgxpool.Pool) CycleRepository {
	return &cycleRepository{pool: pool}
}

const cycleCols = `c.id, c.plan_id, c.cycle_number, c.label, c.started_at, c.completed_at, c.notes`

func scanCycle(row pgx.Row) (entity.Cycle, error) {
	var c entity.Cycle
	err := row.Scan(&c.ID, &c.PlanID, &c.CycleNumber, &c.Label, &c.StartedAt, &c.CompletedAt, &c.Notes)

	return c, err
}

func (r *cycleRepository) CreateCycle(ctx context.Context, userID, planID uuid.UUID, input domain.CreateCycleInput) (entity.Cycle, error) {
	c, err := scanCycle(r.pool.QueryRow(ctx,
		`INSERT INTO cycles AS c (plan_id, cycle_number, label, notes)
		 SELECT $1, COALESCE((SELECT MAX(cycle_number) FROM cycles WHERE plan_id=$1),0)+1, $2, $3
		 WHERE EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$4)
		 RETURNING `+cycleCols,
		planID, input.Label, input.Notes, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Cycle{}, commonErrors.ErrNotFound
	}

	return c, err
}

func (r *cycleRepository) ListCycles(ctx context.Context, userID, planID uuid.UUID) ([]entity.Cycle, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+cycleCols+` FROM cycles c
		 WHERE c.plan_id=$1 AND EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$2)
		 ORDER BY c.cycle_number DESC`, planID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []entity.Cycle{}
	for rows.Next() {
		c, err := scanCycle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}

	return out, rows.Err()
}

func (r *cycleRepository) UpdateCycle(ctx context.Context, userID, id uuid.UUID, input domain.UpdateCycleInput) (entity.Cycle, error) {
	c, err := scanCycle(r.pool.QueryRow(ctx,
		`UPDATE cycles c SET label=$3, notes=$4, completed_at=$5
		 FROM plans p
		 WHERE c.id=$1 AND c.plan_id=p.id AND p.user_id=$2
		 RETURNING `+cycleCols,
		id, userID, input.Label, input.Notes, input.CompletedAt))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Cycle{}, commonErrors.ErrNotFound
	}

	return c, err
}

func (r *cycleRepository) DeleteCycle(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM cycles c USING plans p
		 WHERE c.id=$1 AND c.plan_id=p.id AND p.user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return commonErrors.ErrNotFound
	}

	return nil
}
