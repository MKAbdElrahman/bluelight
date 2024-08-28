package handlers

import (
	"net/http"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	errorhandler "bluelight.mkcodedev.com/src/api/handlers/errorhandler"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newListMovieHandlerFunc(em *errorhandler.ErrorHandeler, movieService *domain.MovieService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Request
		req, requestErr := v1.NewListMoviesRequest(r)
		if requestErr != nil {
			em.SendClientError(w, r, requestErr)
			return
		}
		movies, paginationMetadata, err := movieService.GetAllMovies(domain.MovieFilters{
			Title:    req.QueryParams.Title,
			Genres:   req.QueryParams.Genres,
			Page:     req.QueryParams.Page,
			PageSize: req.QueryParams.PageSize,
			Sort:     req.QueryParams.Sort,
		})

		if err != nil {
			em.SendServerError(w, r, v1.ServerError{
				Code:            http.StatusInternalServerError,
				InternalMessage: err.Error(),
			})
			return
		}

		res := v1.NewListMoviesResponse(movies, paginationMetadata)

		err = jsonio.SendJSON(w, res, res.Status(), res.Headers())
		if err != nil {
			em.SendServerError(w, r, v1.ServerError{
				Code:            http.StatusInternalServerError,
				InternalMessage: err.Error(),
			})
			return
		}
	}
}
