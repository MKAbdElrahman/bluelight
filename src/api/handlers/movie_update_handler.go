package handlers

import (
	"errors"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newUpdateMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *domain.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewUpdateMovieRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		// Business

		m, domainErr := movieService.GetMovie(req.IdPathParam)
		if domainErr != nil {
			switch {
			case errors.Is(domainErr, domain.ErrRecordNotFound):
				em.SendClientError(w, r, v1.NotFoundError)
			default:
				em.SendServerError(w, r, v1.InternalServerError)
			}
			return
		}

		
		if req.Body.Title != nil {
			m.Title = *req.Body.Title
		}
		if req.Body.Year != nil {
			m.Year = *req.Body.Year
		}
		if req.Body.Runtime != nil {
			m.RuntimeInMinutes = *req.Body.Runtime
		}
		if req.Body.Genres != nil {
			m.Genres = req.Body.Genres
		}

		
		domainErr = movieService.UpdateMovie(m)
		if domainErr != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}

		// Response
		res := v1.UpdateMovieResponse{
			Id:               m.Id,
			Title:            m.Title,
			Year:             m.Year,
			Version:          m.Version,
			RuntimeInMinutes: m.RuntimeInMinutes,
			Genres:           m.Genres,
		}

		err := jsonio.SendJSON(w, jsonio.Envelope{"movie": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}
	}
}
