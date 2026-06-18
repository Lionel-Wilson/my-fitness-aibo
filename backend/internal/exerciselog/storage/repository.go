package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	"github.com/lionel/my-fitness-aibo/backend/internal/exerciselog/domain"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

type ExerciseLogRepository interface {
	UpsertExerciseLog(ctx context.Context, userID uuid.UUID, input domain.UpsertExerciseLogInput) (entity.ExerciseLog, error)
	ListExerciseLogs(ctx context.Context, userID, exerciseID uuid.UUID) ([]entity.ExerciseLog, error)
}

type exerciseLogRepository struct {
	pool *pgxpool.Pool
}

func NewExerciseLogRepository(pool *pgxpool.Pool) ExerciseLogRepository {
	return &exerciseLogRepository{pool: pool}
}

func (r *exerciseLogRepository) UpsertExerciseLog(ctx context.Context, userID uuid.UUID, input domain.UpsertExerciseLogInput) (entity.ExerciseLog, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return entity.ExerciseLog{}, err
	}
	defer tx.Rollback(ctx)

	var log entity.ExerciseLog
	err = tx.QueryRow(ctx,
		`INSERT INTO exercise_logs (exercise_id, cycle_id, note, working_weight_kg)
		 SELECT $1,$2,$3,$4
		 WHERE EXISTS (SELECT 1 FROM exercises e JOIN workouts w ON w.id=e.workout_id JOIN plans p ON p.id=w.plan_id WHERE e.id=$1 AND p.user_id=$5)
		   AND EXISTS (SELECT 1 FROM cycles c JOIN plans p ON p.id=c.plan_id WHERE c.id=$2 AND p.user_id=$5)
		 ON CONFLICT (exercise_id, cycle_id)
		   DO UPDATE SET note=EXCLUDED.note, working_weight_kg=EXCLUDED.working_weight_kg
		 RETURNING id, exercise_id, cycle_id, note, working_weight_kg, created_at`,
		input.ExerciseID, input.CycleID, input.Note, input.WorkingWeightKg, userID).
		Scan(&log.ID, &log.ExerciseID, &log.CycleID, &log.Note, &log.WorkingWeightKg, &log.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.ExerciseLog{}, commonErrors.ErrNotFound
	}
	if err != nil {
		return entity.ExerciseLog{}, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM set_logs WHERE exercise_log_id=$1`, log.ID); err != nil {
		return entity.ExerciseLog{}, err
	}

	log.Sets = []entity.SetLog{}
	for _, sp := range input.Sets {
		side := sp.Side
		if side != "left" && side != "right" {
			side = "both"
		}

		var sl entity.SetLog
		err := tx.QueryRow(ctx,
			`INSERT INTO set_logs (exercise_log_id, set_index, side, weight_kg, reps, rpe, is_drop_set)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 RETURNING id, exercise_log_id, set_index, side, weight_kg, reps, rpe, is_drop_set`,
			log.ID, sp.SetIndex, side, sp.WeightKg, sp.Reps, sp.Rpe, sp.IsDropSet).
			Scan(&sl.ID, &sl.ExerciseLogID, &sl.SetIndex, &sl.Side, &sl.WeightKg, &sl.Reps, &sl.Rpe, &sl.IsDropSet)
		if err != nil {
			return entity.ExerciseLog{}, err
		}
		log.Sets = append(log.Sets, sl)
	}

	if err := tx.Commit(ctx); err != nil {
		return entity.ExerciseLog{}, err
	}

	return log, nil
}

func (r *exerciseLogRepository) ListExerciseLogs(ctx context.Context, userID, exerciseID uuid.UUID) ([]entity.ExerciseLog, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT el.id, el.exercise_id, el.cycle_id, el.note, el.working_weight_kg, el.created_at
		 FROM exercise_logs el
		 JOIN cycles c ON c.id = el.cycle_id
		 JOIN exercises e ON e.id = el.exercise_id
		 JOIN workouts w ON w.id = e.workout_id
		 JOIN plans p ON p.id = w.plan_id
		 WHERE el.exercise_id=$1 AND p.user_id=$2
		 ORDER BY c.cycle_number`, exerciseID, userID)
	if err != nil {
		return nil, err
	}

	logs := []entity.ExerciseLog{}
	byID := map[uuid.UUID]int{}
	for rows.Next() {
		var l entity.ExerciseLog
		if err := rows.Scan(&l.ID, &l.ExerciseID, &l.CycleID, &l.Note, &l.WorkingWeightKg, &l.CreatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		l.Sets = []entity.SetLog{}
		byID[l.ID] = len(logs)
		logs = append(logs, l)
	}
	rows.Close()
	if err := rows.Err(); err != nil {
		return nil, err
	}
	if len(logs) == 0 {
		return logs, nil
	}

	ids := make([]uuid.UUID, 0, len(logs))
	for _, l := range logs {
		ids = append(ids, l.ID)
	}

	setRows, err := r.pool.Query(ctx,
		`SELECT id, exercise_log_id, set_index, side, weight_kg, reps, rpe, is_drop_set
		 FROM set_logs WHERE exercise_log_id = ANY($1) ORDER BY set_index, side`, ids)
	if err != nil {
		return nil, err
	}
	defer setRows.Close()

	for setRows.Next() {
		var sl entity.SetLog
		if err := setRows.Scan(&sl.ID, &sl.ExerciseLogID, &sl.SetIndex, &sl.Side, &sl.WeightKg, &sl.Reps, &sl.Rpe, &sl.IsDropSet); err != nil {
			return nil, err
		}
		if idx, ok := byID[sl.ExerciseLogID]; ok {
			logs[idx].Sets = append(logs[idx].Sets, sl)
		}
	}

	return logs, setRows.Err()
}
