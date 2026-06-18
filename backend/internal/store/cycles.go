package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Cycle is the plan-level progression unit (C1, C2, …). One cycle is completing
// each workout of the plan once.
type Cycle struct {
	ID          uuid.UUID  `json:"id"`
	PlanID      uuid.UUID  `json:"planId"`
	CycleNumber int        `json:"cycleNumber"`
	Label       string     `json:"label"`
	StartedAt   time.Time  `json:"startedAt"`
	CompletedAt *time.Time `json:"completedAt"`
	Notes       string     `json:"notes"`
}

const cycleCols = `c.id, c.plan_id, c.cycle_number, c.label, c.started_at, c.completed_at, c.notes`

func scanCycle(row pgx.Row) (Cycle, error) {
	var c Cycle
	err := row.Scan(&c.ID, &c.PlanID, &c.CycleNumber, &c.Label, &c.StartedAt, &c.CompletedAt, &c.Notes)
	return c, err
}

// CreateCycle starts the next cycle for a plan owned by the user. The cycle
// number is automatically the current maximum plus one.
func (s *Store) CreateCycle(ctx context.Context, userID, planID uuid.UUID, label, notes string) (Cycle, error) {
	c, err := scanCycle(s.pool.QueryRow(ctx,
		`INSERT INTO cycles AS c (plan_id, cycle_number, label, notes)
		 SELECT $1, COALESCE((SELECT MAX(cycle_number) FROM cycles WHERE plan_id=$1),0)+1, $2, $3
		 WHERE EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$4)
		 RETURNING `+cycleCols,
		planID, label, notes, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Cycle{}, ErrNotFound
	}
	return c, err
}

// ListCycles returns the cycles of a plan owned by the user, newest first.
func (s *Store) ListCycles(ctx context.Context, userID, planID uuid.UUID) ([]Cycle, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+cycleCols+` FROM cycles c
		 WHERE c.plan_id=$1 AND EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$2)
		 ORDER BY c.cycle_number DESC`, planID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Cycle{}
	for rows.Next() {
		c, err := scanCycle(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, c)
	}
	return out, rows.Err()
}

// UpdateCycle updates the label, notes and completion of a cycle owned by the user.
func (s *Store) UpdateCycle(ctx context.Context, userID, id uuid.UUID, label, notes string, completedAt *time.Time) (Cycle, error) {
	c, err := scanCycle(s.pool.QueryRow(ctx,
		`UPDATE cycles c SET label=$3, notes=$4, completed_at=$5
		 FROM plans p
		 WHERE c.id=$1 AND c.plan_id=p.id AND p.user_id=$2
		 RETURNING `+cycleCols,
		id, userID, label, notes, completedAt))
	if errors.Is(err, pgx.ErrNoRows) {
		return Cycle{}, ErrNotFound
	}
	return c, err
}

// DeleteCycle removes a cycle owned by the user.
func (s *Store) DeleteCycle(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM cycles c USING plans p
		 WHERE c.id=$1 AND c.plan_id=p.id AND p.user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
