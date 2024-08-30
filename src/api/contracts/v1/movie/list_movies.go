package v1movie

import (
	"net/http"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/contracts/webutil"
	"bluelight.mkcodedev.com/src/core/domain/movie"
	"bluelight.mkcodedev.com/src/core/domain/verrors"
)

type ListMoviesRequest struct {
	QueryParams ListMoviesRequestQueryParams
}

type ListMoviesRequestQueryParams struct {
	Title    string
	Genres   []string
	Page     int
	PageSize int
	Sort     string
}

func NewListMoviesRequest(r *http.Request) (ListMoviesRequest, *apierror.ClientError) {
	req := ListMoviesRequest{}

	qParams, errDetails := newListMoviesRequestQueryParams(r)
	if errDetails != nil {
		return ListMoviesRequest{}, apierror.BadRequestError.WithValidationError(errDetails)
	}

	req.QueryParams = qParams

	return req, nil
}

func newListMoviesRequestQueryParams(r *http.Request) (ListMoviesRequestQueryParams, *verrors.ValidationError) {
	qParams := ListMoviesRequestQueryParams{
		Title:    webutil.GetQueryParam(r, "title"),
		Genres:   webutil.GetQueryParamSlice(r, "genres", ","),
		Page:     1,
		PageSize: 20,
		Sort:     "id",
	}

	// Parse and validate the "page" query parameter
	page, err := webutil.GetQueryParamInt(r, "page")

	if err != nil {
		return ListMoviesRequestQueryParams{}, verrors.NewValidationError("page", "invalid page number")
	}

	if page != 0 {
		qParams.Page = page
	}

	// Parse and validate the "pageSize" query parameter
	pageSize, err := webutil.GetQueryParamInt(r, "page_size")
	if err != nil {
		return ListMoviesRequestQueryParams{}, verrors.NewValidationError("page_size", "invalid page size")
	}
	if pageSize != 0 {
		qParams.PageSize = pageSize
	}

	if sort := webutil.GetQueryParam(r, "sort"); sort != "" {
		qParams.Sort = sort
	}

	return qParams, nil
}

type ListMoviesResponse struct {
	Movies             []ShowMovieResponse                `json:"movies"`
	MoviesListMetaData movie.MoviesListPaginationMetadata `json:"pagination_metadata"`
}

func NewListMoviesResponse(movies []*movie.Movie, metadata movie.MoviesListPaginationMetadata) ListMoviesResponse {

	var moviesList ListMoviesResponse
	moviesList.MoviesListMetaData = metadata
	for _, m := range movies {
		moviesList.Movies = append(moviesList.Movies, ShowMovieResponse{
			Id:               m.Id,
			Title:            m.Title,
			Genres:           m.Genres,
			RuntimeInMinutes: m.RuntimeInMinutes,
			Year:             m.Year,
			Version:          m.Version,
		})
	}
	return moviesList
}

func (r ListMoviesResponse) Status() int {
	return http.StatusOK
}
func (r ListMoviesResponse) Headers() http.Header {
	return make(http.Header)
}
