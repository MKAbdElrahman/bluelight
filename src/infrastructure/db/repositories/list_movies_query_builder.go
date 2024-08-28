package repositories

import (
	"fmt"
	"strings"
)

type selectQueryBuilder struct {
	baseQuery  string
	conditions []string
	sort       string
	pageSize   int
	page       int
	argIndex   int
	args       []interface{}
}

func newSelectQueryBuilder() *selectQueryBuilder {
	return &selectQueryBuilder{
		baseQuery: `SELECT id, created_at, title, year, runtime, genres, version FROM movies`,
	}
}

func (b *selectQueryBuilder) setPagination(page int, pageSize int) *selectQueryBuilder {
	b.page = page
	b.pageSize = pageSize
	return b
}
func (b *selectQueryBuilder) build() (string, string, []interface{}, []interface{}) {
	query := b.baseQuery
	countQuery := "SELECT COUNT(*) FROM movies"

	// Apply conditions to both the main query and the count query
	if len(b.conditions) > 0 {
		conditionClause := " WHERE " + strings.Join(b.conditions, " AND ")
		query += conditionClause
		countQuery += conditionClause
	}

	// Sorting (only for the main query)
	switch b.sort {
	case "id", "title", "year", "runtime":
		query += " ORDER BY " + b.sort + " ASC"
	case "-id", "-title", "-year", "-runtime":
		query += " ORDER BY " + b.sort[1:] + " DESC"
	default:
		query += " ORDER BY id ASC" // Default sort by id
	}

	// Arguments for the count query (only the conditions)
	countArgs := make([]interface{}, len(b.args))
	copy(countArgs, b.args)

	// Pagination (only for the main query)
	if b.pageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", b.argIndex+1, b.argIndex+2)
		b.args = append(b.args, b.pageSize, (b.page-1)*b.pageSize)
	}

	// Return the main query, count query, and the arguments for each
	return query, countQuery, b.args, countArgs
}

func (b *selectQueryBuilder) setSort(sort string) *selectQueryBuilder {
	b.sort = sort
	return b
}

func (b *selectQueryBuilder) addTitleFilter(title string) *selectQueryBuilder {
	if title != "" {
		b.conditions = append(b.conditions, fmt.Sprintf("title ILIKE $%d", b.argIndex+1))
		b.args = append(b.args, "%"+title+"%")
		b.argIndex++
	}
	return b
}

func (b *selectQueryBuilder) addGenresFilter(genres []string) *selectQueryBuilder {
	if len(genres) > 0 {
		genrePlaceholders := make([]string, len(genres))
		for i := range genres {
			genrePlaceholders[i] = fmt.Sprintf("$%d", b.argIndex+1)
			b.args = append(b.args, strings.ToLower(genres[i]))
			b.argIndex++
		}
		b.conditions = append(b.conditions, fmt.Sprintf("EXISTS (SELECT 1 FROM unnest(genres) AS genre WHERE LOWER(genre) = ANY (ARRAY[%s]))", strings.Join(genrePlaceholders, ", ")))
	}
	return b
}
