package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// Exercise is a movement within a workout, with target prescription and tips.
type Exercise struct {
	ID          uuid.UUID `json:"id"`
	WorkoutID   uuid.UUID `json:"workoutId"`
	Name        string    `json:"name"`
	OrderIndex  int       `json:"orderIndex"`
	TargetSets  *int      `json:"targetSets"`
	RepLow      *int      `json:"repLow"`
	RepHigh     *int      `json:"repHigh"`
	RpeLow      *float64  `json:"rpeLow"`
	RpeHigh     *float64  `json:"rpeHigh"`
	RestSeconds *int      `json:"restSeconds"`
	Instructions string   `json:"instructions"`
	OrGroup     string    `json:"orGroup"`
	IsOptional  bool      `json:"isOptional"`
	IsUnilateral bool     `json:"isUnilateral"`
	CreatedAt   time.Time `json:"createdAt"`
}

// ExerciseParams holds the editable fields of an exercise.
type ExerciseParams struct {
	Name         string
	OrderIndex   int
	TargetSets   *int
	RepLow       *int
	RepHigh      *int
	RpeLow       *float64
	RpeHigh      *float64
	RestSeconds  *int
	Instructions string
	OrGroup      string
	IsOptional   bool
	IsUnilateral bool
}

const exerciseCols = `e.id, e.workout_id, e.name, e.order_index, e.target_sets, e.rep_low, e.rep_high,
	e.rpe_low, e.rpe_high, e.rest_seconds, e.instructions, e.or_group, e.is_optional, e.is_unilateral, e.created_at`

func scanExercise(row pgx.Row) (Exercise, error) {
	var e Exercise
	err := row.Scan(&e.ID, &e.WorkoutID, &e.Name, &e.OrderIndex, &e.TargetSets,
		&e.RepLow, &e.RepHigh, &e.RpeLow, &e.RpeHigh, &e.RestSeconds,
		&e.Instructions, &e.OrGroup, &e.IsOptional, &e.IsUnilateral, &e.CreatedAt)
	return e, err
}

// CreateExercise inserts an exercise, verifying the workout belongs to the user.
func (s *Store) CreateExercise(ctx context.Context, userID, workoutID uuid.UUID, p ExerciseParams) (Exercise, error) {
	e, err := scanExercise(s.pool.QueryRow(ctx,
		`INSERT INTO exercises AS e (workout_id, name, order_index, target_sets, rep_low, rep_high,
			rpe_low, rpe_high, rest_seconds, instructions, or_group, is_optional, is_unilateral)
		 SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		 WHERE EXISTS (SELECT 1 FROM workouts w JOIN plans p ON p.id=w.plan_id WHERE w.id=$1 AND p.user_id=$14)
		 RETURNING `+exerciseCols,
		workoutID, p.Name, p.OrderIndex, p.TargetSets, p.RepLow, p.RepHigh,
		p.RpeLow, p.RpeHigh, p.RestSeconds, p.Instructions, p.OrGroup, p.IsOptional, p.IsUnilateral, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Exercise{}, ErrNotFound
	}
	return e, err
}

// ListExercises returns the exercises of a workout owned by the user, in order.
func (s *Store) ListExercises(ctx context.Context, userID, workoutID uuid.UUID) ([]Exercise, error) {
	rows, err := s.pool.Query(ctx,
		`SELECT `+exerciseCols+` FROM exercises e
		 WHERE e.workout_id=$1
		   AND EXISTS (SELECT 1 FROM workouts w JOIN plans p ON p.id=w.plan_id WHERE w.id=$1 AND p.user_id=$2)
		 ORDER BY e.order_index, e.created_at`, workoutID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []Exercise{}
	for rows.Next() {
		e, err := scanExercise(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

// GetExercise returns an exercise owned by the user.
func (s *Store) GetExercise(ctx context.Context, userID, id uuid.UUID) (Exercise, error) {
	e, err := scanExercise(s.pool.QueryRow(ctx,
		`SELECT `+exerciseCols+` FROM exercises e
		 JOIN workouts w ON w.id=e.workout_id
		 JOIN plans p ON p.id=w.plan_id
		 WHERE e.id=$1 AND p.user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return Exercise{}, ErrNotFound
	}
	return e, err
}

// UpdateExercise updates an exercise owned by the user.
func (s *Store) UpdateExercise(ctx context.Context, userID, id uuid.UUID, p ExerciseParams) (Exercise, error) {
	e, err := scanExercise(s.pool.QueryRow(ctx,
		`UPDATE exercises e SET name=$3, order_index=$4, target_sets=$5, rep_low=$6, rep_high=$7,
			rpe_low=$8, rpe_high=$9, rest_seconds=$10, instructions=$11, or_group=$12, is_optional=$13, is_unilateral=$14
		 FROM workouts w, plans p
		 WHERE e.id=$1 AND e.workout_id=w.id AND w.plan_id=p.id AND p.user_id=$2
		 RETURNING `+exerciseCols,
		id, userID, p.Name, p.OrderIndex, p.TargetSets, p.RepLow, p.RepHigh,
		p.RpeLow, p.RpeHigh, p.RestSeconds, p.Instructions, p.OrGroup, p.IsOptional, p.IsUnilateral))
	if errors.Is(err, pgx.ErrNoRows) {
		return Exercise{}, ErrNotFound
	}
	return e, err
}

// DeleteExercise removes an exercise owned by the user.
func (s *Store) DeleteExercise(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := s.pool.Exec(ctx,
		`DELETE FROM exercises e USING workouts w, plans p
		 WHERE e.id=$1 AND e.workout_id=w.id AND w.plan_id=p.id AND p.user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
