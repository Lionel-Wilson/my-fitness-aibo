package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Workout is a single training day within a plan.
type Workout struct {
	ID          uuid.UUID `json:"id"`
	PlanID      uuid.UUID `json:"planId"`
	Name        string    `json:"name"`
	DayLabel    string    `json:"dayLabel"`
	OrderIndex  int       `json:"orderIndex"`
	DurationMin *int      `json:"durationMin"`
	Notes       string    `json:"notes"`
	CreatedAt   time.Time `json:"createdAt"`
}

// WorkoutParams holds the editable fields of a workout.
type WorkoutParams struct {
	Name        string
	DayLabel    string
	OrderIndex  int
	DurationMin *int
	Notes       string
}

const workoutCols = `w.id, w.plan_id, w.name, w.day_label, w.order_index, w.duration_min, w.notes, w.created_at`

func scanWorkout(row pgx.Row) (Workout, error) {
	var w Workout
	err := row.Scan(&w.ID, &w.PlanID, &w.Name, &w.DayLabel, &w.OrderIndex, &w.DurationMin, &w.Notes, &w.CreatedAt)
	return w, err
}

// CreateWorkout inserts a workout, verifying the plan belongs to the user.
func (s *Store) CreateWorkout(ctx context.Context, userID, planID uuid.UUID, p WorkoutParams) (Workout, error) {
	w, err := scanWorkout(s.pool.QueryRow(ctx,
		`INSERT INTO workouts AS w (plan_id, name, day_label, order_index, duration_min, notes)
		 SELECT $1,$2,$3,$4,$5,$6
		 WHERE EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$7)
		 RETURNING `+workoutCols,
		planID, p.Name, p.DayLabel, p.OrderIndex, p.DurationMin, p.Notes, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Workout{}, ErrNotFound
	}
	return w, err
}

// ListWorkouts returns the workouts of a plan owned by the user, in order.
func (s *Store) ListWorkouts(ctx context.Context, userID, planID uuid.UUID) ([]Workout, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+workoutCols+` FROM workouts w
		 WHERE w.plan_id=$1 AND EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$2)
		 ORDER BY w.order_index, w.created_at`, planID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Workout{}
	for rows.Next() {
		w, err := scanWorkout(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}
	return out, rows.Err()
}

// GetWorkout returns a workout owned by the user (via its plan).
func (s *Store) GetWorkout(ctx context.Context, userID, id uuid.UUID) (Workout, error) {
	w, err := scanWorkout(s.pool.QueryRow(ctx,
		`SELECT `+workoutCols+` FROM workouts w
		 JOIN plans p ON p.id = w.plan_id
		 WHERE w.id=$1 AND p.user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Workout{}, ErrNotFound
	}
	return w, err
}

// UpdateWorkout updates a workout owned by the user.
func (s *Store) UpdateWorkout(ctx context.Context, userID, id uuid.UUID, p WorkoutParams) (Workout, error) {
	w, err := scanWorkout(s.pool.QueryRow(ctx,
		`UPDATE workouts w SET name=$3, day_label=$4, order_index=$5, duration_min=$6, notes=$7
		 FROM plans p
		 WHERE w.id=$1 AND w.plan_id=p.id AND p.user_id=$2
		 RETURNING `+workoutCols,
		id, userID, p.Name, p.DayLabel, p.OrderIndex, p.DurationMin, p.Notes))
	if errors.Is(err, pgx.ErrNoRows) {
		return Workout{}, ErrNotFound
	}
	return w, err
}

// DeleteWorkout removes a workout owned by the user.
func (s *Store) DeleteWorkout(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM workouts w USING plans p
		 WHERE w.id=$1 AND w.plan_id=p.id AND p.user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
