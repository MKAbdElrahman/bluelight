package domain

type MovieService struct {
	movieRepository MovieRepositoty
}

func NewMovieService(r MovieRepositoty) *MovieService {
	return &MovieService{
		movieRepository: r,
	}
}

func (svc *MovieService) CreateMovie(m *Movie) error {
	return svc.movieRepository.Create(m)
}

func (svc *MovieService) GetMovie(id int64) (*Movie, error) {
	return svc.movieRepository.Read(id)
}
