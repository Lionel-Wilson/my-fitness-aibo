CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE users (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email         TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE plans (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id      UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    quality      TEXT NOT NULL DEFAULT '',
    description  TEXT NOT NULL DEFAULT '',
    cycle_label  TEXT NOT NULL DEFAULT 'Cycle',
    period_start DATE,
    period_end   DATE,
    is_active    BOOLEAN NOT NULL DEFAULT true,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX plans_user_id_idx ON plans(user_id);

CREATE TABLE workouts (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id      UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    name         TEXT NOT NULL,
    day_label    TEXT NOT NULL DEFAULT '',
    order_index  INTEGER NOT NULL DEFAULT 0,
    duration_min INTEGER,
    notes        TEXT NOT NULL DEFAULT '',
    created_at   TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX workouts_plan_id_idx ON workouts(plan_id);

CREATE TABLE exercises (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    workout_id    UUID NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    order_index   INTEGER NOT NULL DEFAULT 0,
    target_sets   INTEGER,
    rep_low       INTEGER,
    rep_high      INTEGER,
    rpe_low       NUMERIC(3,1),
    rpe_high      NUMERIC(3,1),
    rest_seconds  INTEGER,
    instructions  TEXT NOT NULL DEFAULT '',
    or_group      TEXT NOT NULL DEFAULT '',
    is_optional   BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now()
);
CREATE INDEX exercises_workout_id_idx ON exercises(workout_id);

CREATE TABLE cycles (
    id           UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    plan_id      UUID NOT NULL REFERENCES plans(id) ON DELETE CASCADE,
    cycle_number INTEGER NOT NULL,
    label        TEXT NOT NULL DEFAULT '',
    started_at   TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at TIMESTAMPTZ,
    notes        TEXT NOT NULL DEFAULT '',
    UNIQUE (plan_id, cycle_number)
);
CREATE INDEX cycles_plan_id_idx ON cycles(plan_id);

CREATE TABLE exercise_logs (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_id       UUID NOT NULL REFERENCES exercises(id) ON DELETE CASCADE,
    cycle_id          UUID NOT NULL REFERENCES cycles(id) ON DELETE CASCADE,
    note              TEXT NOT NULL DEFAULT '',
    working_weight_kg NUMERIC(6,2),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT now(),
    UNIQUE (exercise_id, cycle_id)
);
CREATE INDEX exercise_logs_exercise_id_idx ON exercise_logs(exercise_id);
CREATE INDEX exercise_logs_cycle_id_idx ON exercise_logs(cycle_id);

CREATE TABLE set_logs (
    id              UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    exercise_log_id UUID NOT NULL REFERENCES exercise_logs(id) ON DELETE CASCADE,
    set_index       INTEGER NOT NULL,
    weight_kg       NUMERIC(6,2),
    reps            INTEGER,
    rpe             NUMERIC(3,1),
    is_drop_set     BOOLEAN NOT NULL DEFAULT false,
    UNIQUE (exercise_log_id, set_index)
);
CREATE INDEX set_logs_exercise_log_id_idx ON set_logs(exercise_log_id);
