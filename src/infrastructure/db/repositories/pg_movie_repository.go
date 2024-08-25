package repositories

import (
	"database/sql"

	"bluelight.mkcodedev.com/src/core/domain"
	"github.com/lib/pq"
)

type postgresMovieRepositry struct {
	db *sql.DB
}

func NewPostgresMovieRepository(db *sql.DB) *postgresMovieRepositry {
	return &postgresMovieRepositry{
		db: db,
	}
}

func (r *postgresMovieRepositry) Create(m *domain.Movie) error {
	query := `
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	args := []any{m.Title, m.Year, m.RuntimeInMinutes, pq.Array(m.Genres)}

	return r.db.QueryRow(query, args...).Scan(&m.Id, &m.CreatedAt, &m.Version)
}

func (r *postgresMovieRepositry) Read(id int64) (*domain.Movie, error) {
	return nil, nil
}

func (r *postgresMovieRepositry) Update(m *domain.Movie) error {
	return nil
}

func (r *postgresMovieRepositry) Delete(id int64) error {
	return nil
}
