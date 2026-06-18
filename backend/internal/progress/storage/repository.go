package storage

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
)

type ProgressRepository interface {
	ExerciseProgress(ctx context.Context, userID, exerciseID uuid.UUID) ([]entity.ProgressPoint, error)
}

type progressRepository struct {
	pool *pgxpool.Pool
}

func NewProgressRepository(pool *pgxpool.Pool) ProgressRepository {
	return &progressRepository{pool: pool}
}

func (r *progressRepository) ExerciseProgress(ctx context.Context, userID, exerciseID uuid.UUID) ([]entity.ProgressPoint, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT c.id, c.cycle_number, c.label, sl.side,
		        COALESCE(MAX(sl.weight_kg), 0) AS top_weight,
		        COALESCE(SUM(sl.weight_kg * sl.reps), 0) AS volume,
		        COALESCE(MAX(sl.weight_kg * (1 + sl.reps / 30.0)), 0) AS best_e1rm,
		        COALESCE(SUM(sl.reps), 0) AS total_reps,
		        c.started_at
		 FROM cycles c
		 JOIN exercise_logs el ON el.cycle_id = c.id AND el.exercise_id = $1
		 JOIN set_logs sl ON sl.exercise_log_id = el.id
		 JOIN exercises e ON e.id = el.exercise_id
		 JOIN workouts w ON w.id = e.workout_id
		 JOIN plans p ON p.id = w.plan_id
		 WHERE p.user_id = $2 AND sl.weight_kg IS NOT NULL AND sl.reps IS NOT NULL
		 GROUP BY c.id, c.cycle_number, c.label, sl.side, c.started_at
		 ORDER BY c.cycle_number, sl.side`, exerciseID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := []entity.ProgressPoint{}
	for rows.Next() {
		var p entity.ProgressPoint
		if err := rows.Scan(&p.CycleID, &p.CycleNumber, &p.Label, &p.Side, &p.TopWeightKg,
			&p.VolumeKg, &p.BestE1RM, &p.TotalReps, &p.StartedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}

	return out, rows.Err()
}
