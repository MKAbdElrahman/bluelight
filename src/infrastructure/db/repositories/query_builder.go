package repositories

import (
	"fmt"
)

type selectQueryBuilder struct {
	baseQuery string
	pageSize  int
	page      int
	argIndex  int
	args      []interface{}
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

	if b.pageSize > 0 {
		query += fmt.Sprintf(" LIMIT $%d OFFSET $%d", b.argIndex+1, b.argIndex+2)
		b.args = append(b.args, b.pageSize, (b.page-1)*b.pageSize)
	}
	return query, b.args
}
