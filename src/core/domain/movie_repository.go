package domain

type MovieRepositoty interface {
	Create(m *Movie) error
	Read(id int64) (*Movie, error)
	Update(m *Movie) error
	Delete(id int64) error
}
