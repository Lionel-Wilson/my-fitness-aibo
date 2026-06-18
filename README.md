# Fitness Aibo

A mobile-first PWA for tracking training plans, cycles and strength progress —
a structured replacement for tracking workouts in notes, with charts.

- **Frontend:** Next.js (App Router) + React + Tailwind, installable PWA
- **Backend:** Go REST API (chi + pgx), Dockerised
- **Database:** PostgreSQL
- **Deploy target:** Railway

## Data model

```
User → Plan → Workout → Exercise
                Plan → Cycle (C1, C2, … — plan-level progression unit)
Exercise + Cycle → ExerciseLog (per-cycle note + working weight) → SetLog (weight × reps per set)
```

A **cycle** is one round of the whole plan (configurable label per plan — "Cycle",
"Week", "Block"). Each exercise can be logged once per cycle; reps and weight are
stored per set. Progress charts derive **estimated 1RM** (Epley: `weight × (1 + reps/30)`),
**top-set weight** and **total volume** per cycle.

## Project layout

```
backend/    Go API (cmd/api, internal/*, migrations/)
frontend/   Next.js PWA (src/app, src/components, src/lib)
docker-compose.yml   Local dev: postgres + backend + frontend
```

## Local development

### Option A — Docker (everything at once)

```bash
cp .env.example .env          # adjust JWT_SECRET for anything non-local
docker compose up --build
```

- Frontend: http://localhost:3000
- API:      http://localhost:8080/api/health
- Postgres: localhost:5435 (host), `db:5432` inside the network

The backend runs migrations automatically on startup.

### Option B — Run services directly (faster iteration)

```bash
# 1. Postgres
docker run -d --name aibo-pg -e POSTGRES_USER=fitness -e POSTGRES_PASSWORD=fitness \
  -e POSTGRES_DB=fitness -p 5432:5432 postgres:16-alpine

# 2. Backend (migrations run on start)
cd backend
DATABASE_URL="postgres://fitness:fitness@localhost:5432/fitness?sslmode=disable" \
  JWT_SECRET="dev-secret-change-me" PORT=8080 go run ./cmd/api

# 3. Frontend
cd frontend
npm install
echo 'NEXT_PUBLIC_API_URL=http://localhost:8080' > .env.local
npm run dev
```

### Tests

```bash
cd backend && go test ./...     # auth (hash/JWT) + e1RM math
```

## Environment variables

| Variable              | Used by  | Notes |
|-----------------------|----------|-------|
| `DATABASE_URL`        | backend  | Postgres connection string |
| `JWT_SECRET`          | backend  | **Required.** Long random string in production |
| `JWT_TTL_HOURS`       | backend  | Access-token lifetime (default 720 = 30 days) |
| `PORT`                | backend  | Listen port (default 8080) |
| `CORS_ORIGINS`        | backend  | Comma-separated allowed origins (the frontend URL) |
| `NEXT_PUBLIC_API_URL` | frontend | API base URL — **inlined at build time** |

## Deploying to Railway

Create a Railway project with **three services**.

1. **Postgres** — add the managed Postgres plugin. Railway provides `DATABASE_URL`.

2. **Backend** — "Deploy from repo", root directory `backend` (it has a Dockerfile).
   Set variables:
   - `DATABASE_URL` → reference the Postgres service's variable
   - `JWT_SECRET` → a long random string
   - `CORS_ORIGINS` → the frontend's public URL (e.g. `https://aibo-web.up.railway.app`)
   - `PORT` → Railway sets this automatically; the app reads it.
   Migrations run on boot.

3. **Frontend** — "Deploy from repo", root directory `frontend` (Dockerfile).
   Because `NEXT_PUBLIC_API_URL` is baked in at build time, add it as a
   **build-time** variable so it is passed to the Docker `ARG`:
   - `NEXT_PUBLIC_API_URL` → the backend's public URL (e.g. `https://aibo-api.up.railway.app`)

   After the backend URL is known, set `CORS_ORIGINS` on the backend to the
   frontend URL and redeploy both if needed.

### Install on your iPhone

Open the frontend URL in Safari → Share → **Add to Home Screen**. It launches
fullscreen like a native app.

## Roadmap (not yet built)

Bodyweight tracking, BMR/TDEE, body stats and nutrition (MyFitnessPal-style). The
schema reserves room for these; they are intentionally out of the first build.
