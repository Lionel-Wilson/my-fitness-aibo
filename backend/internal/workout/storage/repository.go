package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/workout/domain"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

type WorkoutRepository interface {
	CreateWorkout(ctx context.Context, userID, planID uuid.UUID, input domain.WorkoutInput) (entity.Workout, error)
	ListWorkouts(ctx context.Context, userID, planID uuid.UUID) ([]entity.Workout, error)
	GetWorkout(ctx context.Context, userID, id uuid.UUID) (entity.Workout, error)
	UpdateWorkout(ctx context.Context, userID, id uuid.UUID, input domain.WorkoutInput) (entity.Workout, error)
	DeleteWorkout(ctx context.Context, userID, id uuid.UUID) error
}

type workoutRepository struct {
	pool *pgxpool.Pool
}

func NewWorkoutRepository(pool *pgxpool.Pool) WorkoutRepository {
	return &workoutRepository{pool: pool}
}

const workoutCols = `w.id, w.plan_id, w.name, w.day_label, w.order_index, w.duration_min, w.notes, w.created_at`

func scanWorkout(row pgx.Row) (entity.Workout, error) {
	var w entity.Workout
	err := row.Scan(&w.ID, &w.PlanID, &w.Name, &w.DayLabel, &w.OrderIndex, &w.DurationMin, &w.Notes, &w.CreatedAt)

	return w, err
}

func (r *workoutRepository) CreateWorkout(ctx context.Context, userID, planID uuid.UUID, input domain.WorkoutInput) (entity.Workout, error) {
	w, err := scanWorkout(r.pool.QueryRow(ctx,
		`INSERT INTO workouts AS w (plan_id, name, day_label, order_index, duration_min, notes)
		 SELECT $1,$2,$3,$4,$5,$6
		 WHERE EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$7)
		 RETURNING `+workoutCols,
		planID, input.Name, input.DayLabel, input.OrderIndex, input.DurationMin, input.Notes, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Workout{}, commonErrors.ErrNotFound
	}

	return w, err
}

func (r *workoutRepository) ListWorkouts(ctx context.Context, userID, planID uuid.UUID) ([]entity.Workout, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT `+workoutCols+` FROM workouts w
		 WHERE w.plan_id=$1 AND EXISTS (SELECT 1 FROM plans WHERE id=$1 AND user_id=$2)
		 ORDER BY w.order_index, w.created_at`, planID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []entity.Workout{}
	for rows.Next() {
		w, err := scanWorkout(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, w)
	}

	return out, rows.Err()
}

func (r *workoutRepository) GetWorkout(ctx context.Context, userID, id uuid.UUID) (entity.Workout, error) {
	w, err := scanWorkout(r.pool.QueryRow(ctx,
		`SELECT `+workoutCols+` FROM workouts w
		 JOIN plans p ON p.id = w.plan_id
		 WHERE w.id=$1 AND p.user_id=$2`, id, userID))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Workout{}, commonErrors.ErrNotFound
	}

	return w, err
}

func (r *workoutRepository) UpdateWorkout(ctx context.Context, userID, id uuid.UUID, input domain.WorkoutInput) (entity.Workout, error) {
	w, err := scanWorkout(r.pool.QueryRow(ctx,
		`UPDATE workouts w SET name=$3, day_label=$4, order_index=$5, duration_min=$6, notes=$7
		 FROM plans p
		 WHERE w.id=$1 AND w.plan_id=p.id AND p.user_id=$2
		 RETURNING `+workoutCols,
		id, userID, input.Name, input.DayLabel, input.OrderIndex, input.DurationMin, input.Notes))
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.Workout{}, commonErrors.ErrNotFound
	}

	return w, err
}

func (r *workoutRepository) DeleteWorkout(ctx context.Context, userID, id uuid.UUID) error {
	tag, err := r.pool.Exec(ctx,
		`DELETE FROM workouts w USING plans p
		 WHERE w.id=$1 AND w.plan_id=p.id AND p.user_id=$2`, id, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return commonErrors.ErrNotFound
	}

	return nil
}
