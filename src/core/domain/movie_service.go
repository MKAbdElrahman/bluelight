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
