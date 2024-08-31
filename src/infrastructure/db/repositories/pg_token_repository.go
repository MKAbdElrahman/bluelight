package repositories

import (
	"context"
	"database/sql"
	"time"

	"bluelight.mkcodedev.com/src/core/domain/user"
)

type PostgresTokenRepositryConfig struct {
	Timeout time.Duration
}

type postgresTokenRepositry struct {
	db     *sql.DB
	config PostgresTokenRepositryConfig
}

func NewPostgresTokenRepository(db *sql.DB, config PostgresTokenRepositryConfig) *postgresTokenRepositry {
	return &postgresTokenRepositry{
		db:     db,
		config: config,
	}
}

func (r *postgresTokenRepositry) Create(t *user.Token) error {
	query := `
	INSERT INTO tokens (hash, user_id, expiry, scope)
	VALUES ($1, $2, $3, $4)`
	args := []any{t.Hash, t.UserId, t.Expiry, t.Scope}
	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, args...)
	return err
}

func (r *postgresTokenRepositry) DeleteAllForUser(scope string, userID int64) error {
	query := `
	DELETE FROM tokens
	WHERE scope = $1 AND user_id = $2`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_, err := r.db.ExecContext(ctx, query, scope, userID)
	return err
}
