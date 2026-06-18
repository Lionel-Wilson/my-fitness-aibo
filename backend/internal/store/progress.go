package store

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ProgressPoint is one cycle's aggregated performance for an exercise, used for
// charting strength/volume over time.
type ProgressPoint struct {
	CycleID     uuid.UUID `json:"cycleId"`
	CycleNumber int       `json:"cycleNumber"`
	Label       string    `json:"label"`
	Side        string    `json:"side"`
	TopWeightKg float64   `json:"topWeightKg"`
	VolumeKg    float64   `json:"volumeKg"`
	BestE1RM    float64   `json:"bestE1rm"`
	TotalReps   float64   `json:"totalReps"`
	StartedAt   time.Time `json:"startedAt"`
}

// ExerciseProgress returns per-cycle aggregates (top set weight, total volume,
// best estimated 1RM via Epley) for an exercise owned by the user, ordered by
// cycle number. Cycles with no weighted sets are omitted.
func (s *Store) ExerciseProgress(ctx context.Context, userID, exerciseID uuid.UUID) ([]ProgressPoint, error) {
	// Grouped per cycle AND side. Bilateral exercises yield one "both" row per
	// cycle; unilateral exercises yield separate "left"/"right" rows so each side
	// can be charted independently.
	rows, err := s.pool.Query(ctx,
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

	out := []ProgressPoint{}
	for rows.Next() {
		var p ProgressPoint
		if err := rows.Scan(&p.CycleID, &p.CycleNumber, &p.Label, &p.Side, &p.TopWeightKg,
			&p.VolumeKg, &p.BestE1RM, &p.TotalReps, &p.StartedAt); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}
