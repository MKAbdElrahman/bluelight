package moviehandlers

import (
	"errors"
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/movie"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

func NewListMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *movie.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewListMoviesRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}
		movies, paginationMetadata, err := movieService.ListMovies(movie.MovieFilters{
			Title:    req.QueryParams.Title,
			Genres:   req.QueryParams.Genres,
			Page:     req.QueryParams.Page,
			PageSize: req.QueryParams.PageSize,
			Sort:     req.QueryParams.Sort,
		})

		if err != nil {
			var validErr *verrors.ValidationError
			switch {
			case errors.As(err, &validErr):
				em.SendClientError(w, r, apierror.UnprocessableEntityError.WithValidationError(validErr))
			default:
				em.SendServerError(w, r, apierror.NewInternalServerError(err))
			}
			return
		}

		res := v1.NewListMoviesResponse(movies, paginationMetadata)

		err = jsonio.SendJSON(w, res, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, apierror.NewInternalServerError(err))

			return
		}
	}
}
