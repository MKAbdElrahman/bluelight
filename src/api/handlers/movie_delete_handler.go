package handlers

import (
	"errors"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newDeleteMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *domain.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewDeleteMovieRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		err := movieService.DeleteMovie(req.IdPathParam)

		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				em.SendClientError(w, r, v1.NotFoundError)
			default:
				em.SendServerError(w, r, v1.InternalServerError)
			}
			return
		}

		res := v1.DeleteMovieResponse{}

		err = jsonio.SendJSON(w, jsonio.Envelope{"movie": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}

	}
}
