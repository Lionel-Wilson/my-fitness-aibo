# AI Agent Guidelines — Fitness Aibo Backend

Guidelines for AI agents (Cursor, Claude, etc.) working on this Go API. Follow this architecture for every change.

## Overview

This backend uses **clean architecture / DDD** with three layers:

1. **Handler** (`internal/api/{domain}/`) — HTTP, DTOs, status codes
2. **Service** (`internal/{domain}/`) — business logic
3. **Repository** (`internal/{domain}/storage/`) — database access

**Dependency rule (never violate):**

- Handlers → services only (never repositories)
- Services → repositories and other services
- Repositories → database (`pgxpool`) only

## Directory layout

```
backend/
  cmd/api/main.go              # wire repos → services → handlers → router
  internal/
    api/{domain}/handler.go
    api/{domain}/dto/{request.go,response.go}
    api/{domain}/dto/mapper/
    {domain}/service.go
    {domain}/domain/domain.go
    {domain}/mapper/entity_to_domain.go
    {domain}/storage/repository.go
    entity/                      # DB row structs (no json tags)
    http/router/router.go
    config/ db/ middleware/ metrics/
  pkg/commonlibrary/
    render/ request/ errors/ mappers/ messages/ context/ auth/
  migrations/
```

**Existing domains:** `user`, `auth`, `plan`, `workout`, `exercise`, `cycle`, `exerciselog`, `progress`.

## Model types

| Layer      | Package                         | Purpose                          |
|------------|---------------------------------|----------------------------------|
| Entity     | `internal/entity`               | DB rows; used in repositories    |
| Domain     | `internal/{domain}/domain`      | Business models; used in services |
| DTO        | `internal/api/{domain}/dto`     | JSON request/response; used in handlers |

Map at boundaries only:

- Handler: `dto` ↔ `domain` via `internal/api/{domain}/dto/mapper`
- Service: `domain` ↔ `entity` via `internal/{domain}/mapper`

Do not pass entities to handlers or DTOs to services.

## Layer patterns

### Handler

- Export a `Handler` interface; unexported `handler` struct holds `*log.Logger` and the service.
- Each endpoint is a method returning `http.HandlerFunc`.
- Decode with `request.DecodeAndValidate` (requests implement `Validate() error`).
- Respond with `render.Json`, `render.Error`, or `render.NoContent`.
- Map service errors via `render.HandleServiceErrorResponse` and `render.MapStoreError`.
- Get user ID from `commoncontext.UserID(ctx)` on authenticated routes.

```go
func (h *handler) CreatePlan() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        ctx := r.Context()
        userID, _ := commoncontext.UserID(ctx)

        var req dto.PlanRequest
        if err := request.DecodeAndValidate(r.Body, &req); err != nil {
            render.Error(w, http.StatusBadRequest, messages.InvalidJSONMsg)
            return
        }

        input, err := mapper.PlanRequestToDomain(req)
        if err != nil {
            render.Error(w, http.StatusBadRequest, "dates must be in YYYY-MM-DD format")
            return
        }

        plan, err := h.planService.CreatePlan(ctx, userID, input)
        if err != nil {
            render.HandleServiceErrorResponse(h.logger, w, r, "CreatePlan", err, render.MapStoreError)
            return
        }

        render.Json(w, http.StatusCreated, dto.PlanFromDomain(plan))
    }
}
```

### Service

- Export a `Service` interface; unexported `service` struct holds `*log.Logger` and the repository.
- Method signatures use `context.Context` and domain models.
- Map entity ↔ domain in the service; return sentinel errors from `pkg/commonlibrary/errors`.

### Repository

- Export a `{Domain}Repository` interface; unexported struct holds `*pgxpool.Pool`.
- Return `entity` types; translate `pgx.ErrNoRows` to `commonErrors.ErrNotFound`.
- Keep SQL in `storage/repository.go`; do not add a separate data-access package.

## Shared library (`pkg/commonlibrary`)

Use these instead of ad-hoc helpers:

| Package    | Use for                                      |
|------------|----------------------------------------------|
| `render`   | JSON responses, error mapping                |
| `request`  | Body decode/validate, `PathUUID`             |
| `errors`   | `ErrNotFound`, `ErrConflict`                 |
| `context`  | `UserID`, `WithUserID`                       |
| `auth`     | JWT `TokenManager`, password hash/check      |
| `messages` | Shared user-facing strings                   |
| `mappers`  | `ToSimpleErrorResponse`                      |

## Adding a new feature

Work **bottom-up**:

1. Add migration in `migrations/` if schema changes (keep golang-migrate format).
2. Add entity struct in `internal/entity/`.
3. Add `internal/{domain}/storage/repository.go` (interface + pgx impl).
4. Add `internal/{domain}/domain/domain.go` and `internal/{domain}/mapper/`.
5. Add `internal/{domain}/service.go`.
6. Add `internal/api/{domain}/` (dto, mappers, handler).
7. Register route in `internal/http/router/router.go`.
8. Wire in `cmd/api/main.go` (repo → service → handler).

## What not to do

- Do not add a monolithic `handlers` or `store` package.
- Do not call repositories from handlers.
- Do not put business logic in handlers or repositories.
- Do not use SQLBoiler or switch away from `pgx` unless explicitly asked.
- Do not change API routes or JSON field names without an explicit request (frontend depends on them).
- Do not add zap/viper/mockgen unless explicitly asked — use stdlib `log` and env config.

## Wiring checklist (`cmd/api/main.go`)

```go
repo := storage.NewXRepository(pool)
svc := x.NewService(logger, repo)
handler := apix.NewHandler(logger, svc)
// add to router.Handlers{...}
```

## Tests

```bash
cd backend && go build ./... && go vet ./... && go test ./...
```

Run after every backend change.
