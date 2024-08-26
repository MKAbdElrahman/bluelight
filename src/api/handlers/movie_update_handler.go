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

		parsedId, err := parseIdFromPath(r)
		if err != nil || parsedId < 1 {
			em.SendClientError(w, r, v1.NotFoundError)
			return
		}

		m, err := movieService.GetMovie(int64(parsedId))

		if err != nil {
			switch {
			case errors.Is(err, domain.ErrRecordNotFound):
				em.SendClientError(w, r, v1.NotFoundError)
			default:
				em.SendServerError(w, r, v1.InternalServerError)
			}
			return
		}

		var request v1.UpdateMovieRequest
		err = jsonio.NewJSONReader().ReadJSON(r, &request.Body)
		if err != nil {
			em.SendClientError(w, r, v1.BadRequestError.WithDetails(map[string]string{
				"error": err.Error(),
			}))
			return
		}
		vErrors := request.Validate()

		if vErrors.Length() > 0 {
			em.SendClientError(w, r, v1.UnprocessableEntityError.WithDetails(vErrors.Details()))
			return
		}

		m.Title = request.Body.Title
		m.Genres = request.Body.Genres
		m.RuntimeInMinutes = request.Body.Runtime
		m.Year = request.Body.Year

		err = movieService.UpdateMovie(m)
		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}

		jsonio.SendJSON(w, jsonio.Envelope{"movie": v1.UpdateMovieResponse{
			Id:               m.Id,
			Title:            m.Title,
			Year:             m.Year,
			Version:          m.Version,
			RuntimeInMinutes: m.RuntimeInMinutes,
			Genres:           m.Genres,
		}}, http.StatusCreated, nil)
	}
}
