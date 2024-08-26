package handlers

import (
	"fmt"
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newCreateMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *domain.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request v1.CreateMovieRequest
		err := jsonio.NewJSONReader().ReadJSON(r, &request.Body)
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

		m := &domain.Movie{
			Title:            request.Body.Title,
			Year:             int32(request.Body.Year),
			Genres:           request.Body.Genres,
			RuntimeInMinutes: int32(request.Body.Runtime),
		}
		err = movieService.CreateMovie(m)

		if err != nil {
			em.SendServerError(w, r, v1.InternalServerError)
			return
		}

		headers := make(http.Header)
		headers.Set("Location", fmt.Sprintf("/v1/movies/%d", m.Id))

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
