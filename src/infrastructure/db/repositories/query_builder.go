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

func (b *selectQueryBuilder) build() (string, []interface{}) {
	query := b.baseQuery

	if len(b.conditions) > 0 {
		query += " WHERE " + strings.Join(b.conditions, " AND ")
	}

	// Sorting
	switch b.sort {
	case "id", "title", "year", "runtime":
		query += " ORDER BY " + b.sort + " ASC"
	case "-id", "-title", "-year", "-runtime":
		query += " ORDER BY " + b.sort[1:] + " DESC"
	default:
		query += " ORDER BY id ASC" // Default sort by id
	}

	// Pagination
	if b.pageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", b.argIndex+1, b.argIndex+2)
		b.args = append(b.args, b.pageSize, (b.page-1)*b.pageSize)
	}
	return query, b.args
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
