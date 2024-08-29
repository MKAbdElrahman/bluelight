package movie

type MovieService struct {
	movieRepository MovieRepositoty
}

func NewMovieService(r MovieRepositoty) *MovieService {
	return &MovieService{
		movieRepository: r,
	}
}

type MovieCreateParams struct {
	Title            string
	Year             int32
	RuntimeInMinutes int32
	Genres           []string
}

func (svc *MovieService) CreateMovie(params MovieCreateParams) (*Movie, error) {
	m, err := NewMovie(params.Title, params.Year, params.RuntimeInMinutes, params.Genres)
	if err != nil {
		return nil, err
	}

	err = svc.movieRepository.Create(m) // m is mutated
	if err != nil {
		return nil, err
	}

	return m, nil
}

func (svc *MovieService) UpdateMovie(m *Movie) error {
	return svc.movieRepository.Update(m)
}

func (svc *MovieService) GetMovie(id int64) (*Movie, error) {
	return svc.movieRepository.Read(id)
}

func (svc *MovieService) GetAllMovies(filters MovieFilters) ([]*Movie, MoviesListPaginationMetadata, error) {
	return svc.movieRepository.ReadAll(filters)
}

func (svc *MovieService) DeleteMovie(id int64) error {
	return svc.movieRepository.Delete(id)
}
