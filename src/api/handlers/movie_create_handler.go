package handlers

import (
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newCreateMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *domain.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewCreateMovieRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		// Business
		m := &domain.Movie{
			Title:            req.Body.Title,
			Year:             req.Body.Year,
			Genres:           req.Body.Genres,
			RuntimeInMinutes: req.Body.Runtime,
		}

		err := movieService.CreateMovie(m)
		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}

		res := v1.CreateMovieResponse{
			Id:               m.Id,
			Title:            m.Title,
			Year:             m.Year,
			Version:          m.Version,
			RuntimeInMinutes: m.RuntimeInMinutes,
			Genres:           m.Genres,
		}

		err = jsonio.SendJSON(w, jsonio.Envelope{"movie": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}
	}
}
