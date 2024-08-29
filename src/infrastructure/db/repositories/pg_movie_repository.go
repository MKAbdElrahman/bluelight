package repositories

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"github.com/lib/pq"
)

type PostgresMovieRepositryConfig struct {
	Timeout time.Duration
}

type postgresMovieRepositry struct {
	db     *sql.DB
	config PostgresMovieRepositryConfig
}

func NewPostgresMovieRepository(db *sql.DB, config PostgresMovieRepositryConfig) *postgresMovieRepositry {
	return &postgresMovieRepositry{
		db:     db,
		config: config,
	}
}

func (r *postgresMovieRepositry) Create(m *movie.Movie) error {
	query := `
INSERT INTO movies (title, year, runtime, genres)
VALUES ($1, $2, $3, $4)
RETURNING id, created_at, version`
	args := []any{m.Title, m.Year, m.RuntimeInMinutes, pq.Array(m.Genres)}
	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()
	return r.db.QueryRowContext(ctx, query, args...).Scan(&m.Id, &m.CreatedAt, &m.Version)
}

func (r *postgresMovieRepositry) Read(id int64) (*movie.Movie, error) {

	if id < 1 {
		return nil, domain.ErrRecordNotFound
	}

	query := `
	SELECT id, created_at, title, year, runtime, genres, version
	FROM movies
	WHERE id = $1`

	var movie movie.Movie

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, id).Scan(
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

func (r *postgresMovieRepositry) Update(m *movie.Movie) error {
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

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	err := r.db.QueryRowContext(ctx, query, args...).Scan(&m.Version)

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

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	result, err := r.db.ExecContext(ctx, query, id)
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

func (r *postgresMovieRepositry) ReadAll(filters movie.MovieFilters) ([]*movie.Movie, movie.MoviesListPaginationMetadata, error) {
	// Create the query builder with the necessary filters and pagination
	builder := newSelectQueryBuilder().
		setPagination(filters.Page, filters.PageSize).
		setSort(filters.Sort).
		addTitleFilter(filters.Title).
		addGenresFilter(filters.Genres)

	// Build the main query and the count query along with their respective arguments
	query, countQuery, args, countArgs := builder.build()

	ctx, cancel := context.WithTimeout(context.Background(), r.config.Timeout)
	defer cancel()

	// Execute the count query
	var totalCount int
	err := r.db.QueryRowContext(ctx, countQuery, countArgs...).Scan(&totalCount)
	if err != nil {
		return nil, movie.MoviesListPaginationMetadata{}, err
	}

	// Execute the main query
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, movie.MoviesListPaginationMetadata{}, err
	}
	defer rows.Close()

	movies := []*movie.Movie{}

	for rows.Next() {
		var m movie.Movie
		err := rows.Scan(
			&m.Id,
			&m.CreatedAt,
			&m.Title,
			&m.Year,
			&m.RuntimeInMinutes,
			pq.Array(&m.Genres),
			&m.Version,
		)
		if err != nil {
			return nil, movie.MoviesListPaginationMetadata{}, err
		}
		movies = append(movies, &m)
	}
	if err = rows.Err(); err != nil {
		return nil, movie.MoviesListPaginationMetadata{}, err
	}

	// Calculate pagination metadata
	paginationMetadata := movie.MoviesListPaginationMetadata{
		TotalCount: totalCount,
		TotalPages: (totalCount + builder.pageSize - 1) / builder.pageSize,
		Page:       builder.page,
		PageSize:   builder.pageSize,
	}

	return movies, paginationMetadata, nil
}
