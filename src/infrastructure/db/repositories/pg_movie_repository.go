package repositories

import (
	"database/sql"
	"errors"

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

	if id < 1 {
		return nil, domain.ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie domain.Movie

	err := r.db.QueryRow(query, id).Scan(
		&movie.Id,
		&movie.CreatedAt,
		&movie.Title,
		&movie.Year,
		&movie.RuntimeInMinutes,
		pq.Array(&movie.Genres),
		&movie.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, domain.ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (r *postgresMovieRepositry) Update(m *domain.Movie) error {
	query := `
UPDATE movies
SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
WHERE id = $5 AND version = $6
RETURNING version`
	args := []any{
		m.Title,
		m.Year,
		m.RuntimeInMinutes,
		pq.Array(m.Genres),
		m.Id,
		m.Version,
	}

	err := r.db.QueryRow(query, args...).Scan(&m.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return domain.ErrEditConflict
		default:
			return err
		}
	}

	return nil

}

func (r *postgresMovieRepositry) Delete(id int64) error {

	if id < 1 {
		return domain.ErrRecordNotFound
	}
	query := `
	DELETE FROM movies
	WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return domain.ErrRecordNotFound
	}
	return nil
}
