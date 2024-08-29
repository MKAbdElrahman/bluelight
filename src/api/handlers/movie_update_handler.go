package handlers

import (
	"errors"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/movie"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newUpdateMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *movie.MovieService) http.HandlerFunc {
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
			case errors.Is(domainErr, movie.ErrRecordNotFound):
				em.SendClientError(w, r, apierror.NotFoundError)
			default:
				em.SendServerError(w, r, &apierror.ServerError{
					Code:            http.StatusInternalServerError,
					InternalMessage: domainErr.Error(),
				})
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
			switch {
			case errors.Is(domainErr, movie.ErrEditConflict):
				em.SendClientError(w, r, apierror.ConflictError)
			default:
				em.SendServerError(w, r, &apierror.ServerError{
					Code:            http.StatusInternalServerError,
					InternalMessage: domainErr.Error(),
				})
			}
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
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}
	}
}
