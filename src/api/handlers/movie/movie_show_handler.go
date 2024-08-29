package moviehandlers

import (
	"errors"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/movie"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/movie"
)

func NewShowMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *movie.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewShowMovieRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}

		m, err := movieService.GetMovie(req.IdPathParam)

		if err != nil {
			switch {
			case errors.Is(err, movie.ErrRecordNotFound):
				em.SendClientError(w, r, apierror.NotFoundError)
			default:
				em.SendServerError(w, r, &apierror.ServerError{
					Code:            http.StatusInternalServerError,
					InternalMessage: err.Error(),
				})
			}
			return
		}
		res := v1.ShowMovieResponse{
			Id:               m.Id,
			Title:            m.Title,
			Year:             m.Year,
			Version:          m.Version,
			RuntimeInMinutes: m.RuntimeInMinutes,
			Genres:           m.Genres,
		}
		err = jsonio.SendJSON(w, jsonio.Envelope{"movie": res}, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}

	}
}
