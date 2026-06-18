package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/exercise/domain"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

type ExerciseRepository interface {
	CreateExercise(ctx context.Context, userID, workoutID uuid.UUID, input domain.ExerciseInput) (entity.Exercise, error)
	ListExercises(ctx context.Context, userID, workoutID uuid.UUID) ([]entity.Exercise, error)
	GetExercise(ctx context.Context, userID, id uuid.UUID) (entity.Exercise, error)
	UpdateExercise(ctx context.Context, userID, id uuid.UUID, input domain.ExerciseInput) (entity.Exercise, error)
	DeleteExercise(ctx context.Context, userID, id uuid.UUID) error
}

type exerciseRepository struct {
	pool *pgxpool.Pool
}

func NewExerciseRepository(pool *pgxpool.Pool) ExerciseRepository {
	return &exerciseRepository{pool: pool}
}

const exerciseCols = `e.id, e.workout_id, e.name, e.order_index, e.target_sets, e.rep_low, e.rep_high,
	e.rpe_low, e.rpe_high, e.rest_seconds, e.instructions, e.or_group, e.is_optional, e.is_unilateral, e.created_at`

func scanExercise(row pgx.Row) (entity.Exercise, error) {
	var e entity.Exercise
	err := row.Scan(&e.ID, &e.WorkoutID, &e.Name, &e.OrderIndex, &e.TargetSets,
		&e.RepLow, &e.RepHigh, &e.RpeLow, &e.RpeHigh, &e.RestSeconds,
		&e.Instructions, &e.OrGroup, &e.IsOptional, &e.IsUnilateral, &e.CreatedAt)

	return e, err
}

func (r *exerciseRepository) CreateExercise(ctx context.Context, userID, workoutID uuid.UUID, input domain.ExerciseInput) (entity.Exercise, error) {
	e, err := scanExercise(r.pool.QueryRow(ctx,
		`INSERT INTO exercises AS e (workout_id, name, order_index, target_sets, rep_low, rep_high,
			rpe_low, rpe_high, rest_seconds, instructions, or_group, is_optional, is_unilateral)
		 SELECT $1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13
		 WHERE EXISTS (SELECT 1 FROM workouts w JOIN plans p ON p.id=w.plan_id WHERE w.id=$1 AND p.user_id=$14)
		 RETURNING `+exerciseCols,
		workoutID, input.Name, input.OrderIndex, input.TargetSets, input.RepLow, input.RepHigh,
		input.RpeLow, input.RpeHigh, input.RestSeconds, input.Instructions, input.OrGroup, input.IsOptional, input.IsUnilateral, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Exercise{}, commonErrors.ErrNotFound
	}

	return e, err
}

func (r *exerciseRepository) ListExercises(ctx context.Context, userID, workoutID uuid.UUID) ([]entity.Exercise, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+exerciseCols+` FROM exercises e
		 WHERE e.workout_id=$1
		   AND EXISTS (SELECT 1 FROM workouts w JOIN plans p ON p.id=w.plan_id WHERE w.id=$1 AND p.user_id=$2)
		 ORDER BY e.order_index, e.created_at`, workoutID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []entity.Exercise{}
	for rows.Next() {
		e, err := scanExercise(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, e)
	}

	return out, rows.Err()
}

func (r *exerciseRepository) GetExercise(ctx context.Context, userID, id uuid.UUID) (entity.Exercise, error) {
	e, err := scanExercise(r.pool.QueryRow(ctx,
		`SELECT `+exerciseCols+` FROM exercises e
		 JOIN workouts w ON w.id=e.workout_id
		 JOIN plans p ON p.id=w.plan_id
		 WHERE e.id=$1 AND p.user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Exercise{}, commonErrors.ErrNotFound
	}

	return e, err
}

func (r *exerciseRepository) UpdateExercise(ctx context.Context, userID, id uuid.UUID, input domain.ExerciseInput) (entity.Exercise, error) {
	e, err := scanExercise(r.pool.QueryRow(ctx,
		`UPDATE exercises e SET name=$3, order_index=$4, target_sets=$5, rep_low=$6, rep_high=$7,
			rpe_low=$8, rpe_high=$9, rest_seconds=$10, instructions=$11, or_group=$12, is_optional=$13, is_unilateral=$14
		 FROM workouts w, plans p
		 WHERE e.id=$1 AND e.workout_id=w.id AND w.plan_id=p.id AND p.user_id=$2
		 RETURNING `+exerciseCols,
		id, userID, input.Name, input.OrderIndex, input.TargetSets, input.RepLow, input.RepHigh,
		input.RpeLow, input.RpeHigh, input.RestSeconds, input.Instructions, input.OrGroup, input.IsOptional, input.IsUnilateral))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Exercise{}, commonErrors.ErrNotFound
	}

	return e, err
}

func (r *exerciseRepository) DeleteExercise(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM exercises e USING workouts w, plans p
		 WHERE e.id=$1 AND e.workout_id=w.id AND w.plan_id=p.id AND p.user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return commonErrors.ErrNotFound
	}

	return nil
}
