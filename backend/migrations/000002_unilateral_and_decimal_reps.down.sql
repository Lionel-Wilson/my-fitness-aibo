ALTER TABLE set_logs DROP CONSTRAINT set_logs_log_set_side_key;
ALTER TABLE set_logs ADD CONSTRAINT set_logs_exercise_log_id_set_index_key UNIQUE (exercise_log_id, set_index);
ALTER TABLE set_logs DROP CONSTRAINT set_logs_side_check;
ALTER TABLE set_logs DROP COLUMN side;
ALTER TABLE set_logs ALTER COLUMN reps TYPE INTEGER USING ROUND(reps);
ALTER TABLE exercises DROP COLUMN is_unilateral;
