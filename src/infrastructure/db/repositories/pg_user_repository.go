package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"bluelight.mkcodedev.com/src/core/domain/user"
)

type PostgresUserRepositryConfig struct {
	Timeout time.Duration
}

type postgresUserRepositry struct {
	db     *sql.DB
	config PostgresUserRepositryConfig
}

func NewPostgresUserRepository(db *sql.DB, config PostgresUserRepositryConfig) *postgresUserRepositry {
	return &postgresUserRepositry{
		db:     db,
		config: config,
	}
}

func (r *postgresUserRepositry) Create(u *user.User) error {
	query := `
	INSERT INTO users (name, email, password_hash, activated)
	VALUES ($1, $2, $3, $4)
	RETURNING id, created_at, version`

	args := []any{u.Name, u.Email, u.PasswordHash, u.Activated}

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&u.Id, &u.CreatedAt, &u.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return user.ErrDuplicateEmail
		default:
			return err
		}
	}

	return nil
}

func (r *postgresUserRepositry) GetByEmail(email string) (*user.User, error) {

	query := `
	SELECT id, created_at, name, email, password_hash, activated, version
	FROM users
	WHERE email = $1`

	var u user.User

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&u.Id,
		&u.CreatedAt,
		&u.Name,
		&u.Email,
		&u.PasswordHash,
		&u.Activated,
		&u.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, user.ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &u, nil
}

func (r *postgresUserRepositry) Update(u *user.User) error {
	query := `
	UPDATE users
	SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
	WHERE id = $5 AND version = $6
	RETURNING version`
	args := []any{
		u.Name,
		u.Email,
		u.PasswordHash,
		u.Activated,
		u.Id,
		u.Version,
	}
	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&u.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return user.ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return user.ErrEditConflict
		default:
			return err
		}
	}
	return nil
}
