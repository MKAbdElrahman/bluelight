package moviehandlers

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	v1 "bluelight.mkcodedev.com/src/api/contracts/v1/movie"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/api/lib/jsonio"
	"bluelight.mkcodedev.com/src/core/domain/movie"
)

func NewListMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *movie.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewListMoviesRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}
		movies, paginationMetadata, err := movieService.GetAllMovies(movie.MovieFilters{
			Title:    req.QueryParams.Title,
			Genres:   req.QueryParams.Genres,
			Page:     req.QueryParams.Page,
			PageSize: req.QueryParams.PageSize,
			Sort:     req.QueryParams.Sort,
		})

		if err != nil {
			em.SendServerError(w, r, &apierror.ServerError{
				Code:            http.StatusInternalServerError,
				InternalMessage: err.Error(),
			})
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
