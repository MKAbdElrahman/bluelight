package handlers

import (
	"errors"
	"net/http"
	"strconv"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newShowMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *domain.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parsedId, err := parseIdFromPath(r)
		if err != nil || parsedId < 1 {
			em.SendClientError(w, r, v1.NotFoundError)
			return
		}

		m, err := movieService.GetMovie(parsedId)

		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				em.SendClientError(w, r, v1.NotFoundError)
			default:
				em.SendServerError(w, r, v1.InternalServerError)
			}
			return
		}
		err = jsonio.SendJSON(w, jsonio.Envelope{"movie": m}, http.StatusOK, nil)
		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}

	}
}

func parseIdFromPath(r *http.Request) (int64, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.ParseInt(idFromPath, 10, 64)
	return parsedId, err
}
