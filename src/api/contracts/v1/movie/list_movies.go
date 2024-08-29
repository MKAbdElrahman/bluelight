package v1movie

import (
	"net/http"
	"slices"

	"bluelight.mkcodedev.com/src/api/contracts/v1/apierror"
	"bluelight.mkcodedev.com/src/api/contracts/webutil"
	"bluelight.mkcodedev.com/src/core/domain/movie"
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
	if len(errDetails) != 0 {
		return ListMoviesRequest{}, apierror.BadRequestError.WithDetails(errDetails)
	}
	errDetails = qParams.validateRanges()
	if len(errDetails) != 0 {
		return ListMoviesRequest{}, apierror.BadRequestError.WithDetails(errDetails)
	}

	req.QueryParams = qParams

	return req, nil
}

func newListMoviesRequestQueryParams(r *http.Request) (ListMoviesRequestQueryParams, map[string]string) {
	qParams := ListMoviesRequestQueryParams{
		Title:    webutil.GetQueryParam(r, "title"),
		Genres:   webutil.GetQueryParamSlice(r, "genres", ","),
		Page:     1,
		PageSize: 20,
		Sort:     "id",
	}

	errors := make(map[string]string)

	// Parse and validate the "page" query parameter
	if page, err :=  webutil.GetQueryParamInt(r, "page"); err != nil {
		errors["page"] = "invalid page number"
	} else if page != 0 {
		qParams.Page = page
	}

	// Parse and validate the "pageSize" query parameter
	if pageSize, err :=  webutil.GetQueryParamInt(r, "page_size"); err != nil {
		errors["page_size"] = "invalid page size"
	} else if pageSize != 0 {
		qParams.PageSize = pageSize
	}

	// Parse the "sort" query parameter
	if sort :=  webutil.GetQueryParam(r, "sort"); sort != "" {
		qParams.Sort = sort
	}

	return qParams, errors
}

func (q ListMoviesRequestQueryParams) validateRanges() map[string]string {
	errors := make(map[string]string)

	if q.Page <= 0 {
		errors["page"] = "must be greater than zero"
	}

	if q.Page > 10_000_000 {
		errors["page"] = "must be a maximum of 10 million"
	}

	if q.PageSize <= 0 {
		errors["page_size"] = "must be greater than zero"
	}

	if q.PageSize > 10_000_000 {
		errors["page_size"] = "must be a maximum of 10 million"
	}

	if !slices.Contains([]string{"id", "title", "year", "runtime", "-id", "-title", "-year", "-runtime"}, q.Sort) {
		errors["sort"] = "invalid sort value"
	}
	return errors
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
