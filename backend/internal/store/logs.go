package store

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// SetLog is one logged set within an exercise log. Side is "both" for bilateral
// exercises, or "left"/"right" for unilateral ones.
type SetLog struct {
	ID            uuid.UUID `json:"id"`
	ExerciseLogID uuid.UUID `json:"-"`
	SetIndex      int       `json:"setIndex"`
	Side          string    `json:"side"`
	WeightKg      *float64  `json:"weightKg"`
	Reps          *float64  `json:"reps"`
	Rpe           *float64  `json:"rpe"`
	IsDropSet     bool      `json:"isDropSet"`
}

// SetLogParams is the input for a single set.
type SetLogParams struct {
	SetIndex  int      `json:"setIndex"`
	Side      string   `json:"side"`
	WeightKg  *float64 `json:"weightKg"`
	Reps      *float64 `json:"reps"`
	Rpe       *float64 `json:"rpe"`
	IsDropSet bool     `json:"isDropSet"`
}

// ExerciseLog is an exercise's record for one cycle: a note, an optional working
// weight, and the set rows.
type ExerciseLog struct {
	ID              uuid.UUID `json:"id"`
	ExerciseID      uuid.UUID `json:"exerciseId"`
	CycleID         uuid.UUID `json:"cycleId"`
	Note            string    `json:"note"`
	WorkingWeightKg *float64  `json:"workingWeightKg"`
	CreatedAt       time.Time `json:"createdAt"`
	Sets            []SetLog  `json:"sets"`
}

// UpsertExerciseLog creates or replaces the log (note, working weight and all
// sets) for an exercise within a cycle. Both the exercise and the cycle must be
// owned by the user. The set list fully replaces any existing sets.
func (s *Store) UpsertExerciseLog(ctx context.Context, userID, exerciseID, cycleID uuid.UUID, note string, workingWeight *float64, sets []SetLogParams) (ExerciseLog, error) {
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return ExerciseLog{}, err
	}
	defer tx.Rollback(ctx)

	var log ExerciseLog
	err = tx.QueryRow(ctx,
		`INSERT INTO exercise_logs (exercise_id, cycle_id, note, working_weight_kg)
		 SELECT $1,$2,$3,$4
		 WHERE EXISTS (SELECT 1 FROM exercises e JOIN workouts w ON w.id=e.workout_id JOIN plans p ON p.id=w.plan_id WHERE e.id=$1 AND p.user_id=$5)
		   AND EXISTS (SELECT 1 FROM cycles c JOIN plans p ON p.id=c.plan_id WHERE c.id=$2 AND p.user_id=$5)
		 ON CONFLICT (exercise_id, cycle_id)
		   DO UPDATE SET note=EXCLUDED.note, working_weight_kg=EXCLUDED.working_weight_kg
		 RETURNING id, exercise_id, cycle_id, note, working_weight_kg, created_at`,
		exerciseID, cycleID, note, workingWeight, userID).
		Scan(&log.ID, &log.ExerciseID, &log.CycleID, &log.Note, &log.WorkingWeightKg, &log.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return ExerciseLog{}, ErrNotFound
	}
	if err != nil {
		return ExerciseLog{}, err
	}

	if _, err := tx.Exec(ctx, `DELETE FROM set_logs WHERE exercise_log_id=$1`, log.ID); err != nil {
		return ExerciseLog{}, err
	}

	log.Sets = []SetLog{}
	for _, sp := range sets {
		side := sp.Side
		if side != "left" && side != "right" {
			side = "both"
		}
		var sl SetLog
		err := tx.QueryRow(ctx,
			`INSERT INTO set_logs (exercise_log_id, set_index, side, weight_kg, reps, rpe, is_drop_set)
			 VALUES ($1,$2,$3,$4,$5,$6,$7)
			 RETURNING id, exercise_log_id, set_index, side, weight_kg, reps, rpe, is_drop_set`,
			log.ID, sp.SetIndex, side, sp.WeightKg, sp.Reps, sp.Rpe, sp.IsDropSet).
			Scan(&sl.ID, &sl.ExerciseLogID, &sl.SetIndex, &sl.Side, &sl.WeightKg, &sl.Reps, &sl.Rpe, &sl.IsDropSet)
		if err != nil {
			return ExerciseLog{}, err
		}
		log.Sets = append(log.Sets, sl)
	}

	if err := tx.Commit(ctx); err != nil {
		return ExerciseLog{}, err
	}
	return log, nil
}

// ListExerciseLogs returns all logs (with sets) for an exercise owned by the
// user, ordered by cycle number. This powers the cycle history grid.
func (s *Store) ListExerciseLogs(ctx context.Context, userID, exerciseID uuid.UUID) ([]ExerciseLog, error) {
	rows, err := s.pool.Query(ctx,
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

	logs := []ExerciseLog{}
	byID := map[uuid.UUID]int{}
	for rows.Next() {
		var l ExerciseLog
		if err := rows.Scan(&l.ID, &l.ExerciseID, &l.CycleID, &l.Note, &l.WorkingWeightKg, &l.CreatedAt); err != nil {
			rows.Close()
			return nil, err
		}
		l.Sets = []SetLog{}
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
	setRows, err := s.pool.Query(ctx,
		`SELECT id, exercise_log_id, set_index, side, weight_kg, reps, rpe, is_drop_set
		 FROM set_logs WHERE exercise_log_id = ANY($1) ORDER BY set_index, side`, ids)
	if err != nil {
		return nil, err
	}
	defer setRows.Close()
	for setRows.Next() {
		var sl SetLog
		if err := setRows.Scan(&sl.ID, &sl.ExerciseLogID, &sl.SetIndex, &sl.Side, &sl.WeightKg, &sl.Reps, &sl.Rpe, &sl.IsDropSet); err != nil {
			return nil, err
		}
		if idx, ok := byID[sl.ExerciseLogID]; ok {
			logs[idx].Sets = append(logs[idx].Sets, sl)
		}
	}
	return logs, setRows.Err()
}
