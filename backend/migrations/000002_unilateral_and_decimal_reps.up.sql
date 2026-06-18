-- Support unilateral (per-side) exercises and fractional reps.

ALTER TABLE exercises ADD COLUMN is_unilateral BOOLEAN NOT NULL DEFAULT false;

-- Allow fractional reps (e.g. 8.5).
ALTER TABLE set_logs ALTER COLUMN reps TYPE NUMERIC(4,1);

-- Each set can now have per-side rows. 'both' = bilateral (default, unchanged).
ALTER TABLE set_logs ADD COLUMN side TEXT NOT NULL DEFAULT 'both';
ALTER TABLE set_logs ADD CONSTRAINT set_logs_side_check CHECK (side IN ('both', 'left', 'right'));

-- A set is now unique per (log, set_index, side) so left & right can share an index.
ALTER TABLE set_logs DROP CONSTRAINT set_logs_exercise_log_id_set_index_key;
ALTER TABLE set_logs ADD CONSTRAINT set_logs_log_set_side_key UNIQUE (exercise_log_id, set_index, side);
