package moviehandlers

import (
	"errors"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/movie"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/movie"
)

func NewDeleteMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *movie.MovieService) http.HandlerFunc {
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
			case errors.Is(err, movie.ErrRecordNotFound):
				em.SendClientError(w, r, apierror.NotFoundError)
			default:
				em.SendServerError(w, r, apierror.NewInternalServerError(err))
			}
			return
		}

		res := v1.DeleteMovieResponse{}

		err = jsonio.SendJSON(w, jsonio.Envelope{"movie": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}

	}
}
