package movie

import (
	"slices"

	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

type MovieFilters struct {
	Title    string
	Genres   []string
	Page     int
	PageSize int
	Sort     string
}

func (q MovieFilters) validateRanges() error {

	if q.Page <= 0 {
		return verrors.NewValidationError("page", "must be greater than zero")
	}

	if q.Page > 10_000_000 {
		return verrors.NewValidationError("page", "must be a maximum of 10 million")
	}

	if q.PageSize <= 0 {
		return verrors.NewValidationError("page_size", "must be greater than zero")
	}

	if q.PageSize > 10_000_000 {
		return verrors.NewValidationError("page_size", "must be a maximum of 10 million")
	}

	if !slices.Contains([]string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}, q.Sort) {
		return verrors.NewValidationError("sort", "invalid sort value")
	}
	return nil
}
