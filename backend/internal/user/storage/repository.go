package storage

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/lionel/my-fitness-aibo/backend/internal/entity"
	commonErrors "github.com/lionel/my-fitness-aibo/backend/pkg/commonlibrary/errors"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (entity.User, error)
	GetUserByEmail(ctx context.Context, email string) (entity.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error)
}

type userRepository struct {
	pool *pgxpool.Pool
}

func NewUserRepository(pool *pgxpool.Pool) UserRepository {
	return &userRepository{pool: pool}
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash string) (entity.User, error) {
	var u entity.User
	err := r.pool.QueryRow(ctx,
		`INSERT INTO users (email, password_hash)
		 VALUES ($1, $2)
		 RETURNING id, email, password_hash, created_at`,
		email, passwordHash,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return entity.User{}, commonErrors.ErrConflict
		}

		return entity.User{}, err
	}

	return u, nil
}

func (r *userRepository) GetUserByEmail(ctx context.Context, email string) (entity.User, error) {
	var u entity.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE email = $1`,
		email,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, commonErrors.ErrNotFound
	}

	return u, err
}

func (r *userRepository) GetUserByID(ctx context.Context, id uuid.UUID) (entity.User, error) {
	var u entity.User
	err := r.pool.QueryRow(ctx,
		`SELECT id, email, password_hash, created_at FROM users WHERE id = $1`,
		id,
	).Scan(&u.ID, &u.Email, &u.PasswordHash, &u.CreatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return entity.User{}, commonErrors.ErrNotFound
	}

	return u, err
}
