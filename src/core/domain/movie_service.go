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

func (svc *MovieService) UpdateMovie(m *Movie) error {
	return svc.movieRepository.Update(m)
}

func (svc *MovieService) GetMovie(id int64) (*Movie, error) {
	return svc.movieRepository.Read(id)
}

func (svc *MovieService) DeleteMovie(id int64) error {
	return svc.movieRepository.Delete(id)
}
