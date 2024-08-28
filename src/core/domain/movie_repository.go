package domain

type MovieFilters struct {
	Title    string
	Genres   []string
	Page     int
	PageSize int
	Sort     string
}

type MovieRepositoty interface {
	Create(m *Movie) error
	Read(id int64) (*Movie, error)
	ReadAll(filters MovieFilters) ([]*Movie, error)
	Update(m *Movie) error
	Delete(id int64) error
}
