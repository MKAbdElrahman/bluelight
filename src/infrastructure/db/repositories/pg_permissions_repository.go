package repositories

import (
	"context"
	"database/sql"
	"time"

	"bluelight.mkcodedev.com/src/core/domain/user"
)

type PostgresPermissionRepositryConfig struct {
	Timeout time.Duration
}

type postgresPermissionRepositry struct {
	db     *sql.DB
	config PostgresPermissionRepositryConfig
}

func NewPostgresPermissionRepository(db *sql.DB, config PostgresPermissionRepositryConfig) *postgresPermissionRepositry {
	return &postgresPermissionRepositry{
		db:     db,
		config: config,
	}
}

func (r *postgresPermissionRepositry) GetAllForUser(userID int64) (user.Permissions, error) {
	query := `
SELECT permissions.code
FROM permissions
INNER JOIN users_permissions ON users_permissions.permission_id = permissions.id
INNER JOIN users ON users_permissions.user_id = users.id
WHERE users.id = $1`

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var permissions user.Permissions
	for rows.Next() {
		var permission string
		err := rows.Scan(&permission)
		if err != nil {
			return nil, err
		}
		permissions = append(permissions, permission)
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return permissions, nil
}
