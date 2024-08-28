package v1

import (
	"net/http"
	"slices"

	"bluelight.mkcodedev.com/src/core/domain"
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

func NewListMoviesRequest(r *http.Request) (ListMoviesRequest, *ClientError) {
	req := ListMoviesRequest{}

	qParams, errDetails := newListMoviesRequestQueryParams(r)
	if len(errDetails) != 0 {
		return ListMoviesRequest{}, BadRequestError.WithDetails(errDetails)
	}
	errDetails = qParams.validateRanges()
	if len(errDetails) != 0 {
		return ListMoviesRequest{}, BadRequestError.WithDetails(errDetails)
	}

	req.QueryParams = qParams

	return req, nil
}

func newListMoviesRequestQueryParams(r *http.Request) (ListMoviesRequestQueryParams, map[string]string) {
	qParams := ListMoviesRequestQueryParams{
		Title:    GetQueryParam(r, "title"),
		Genres:   GetQueryParamSlice(r, "genres", ","),
		Page:     1,
		PageSize: 20,
		Sort:     "id",
	}

	errors := make(map[string]string)

	// Parse and validate the "page" query parameter
	if page, err := GetQueryParamInt(r, "page"); err != nil {
		errors["page"] = "invalid page number"
	} else if page != 0 {
		qParams.Page = page
	}

	// Parse and validate the "pageSize" query parameter
	if pageSize, err := GetQueryParamInt(r, "page_size"); err != nil {
		errors["page_size"] = "invalid page size"
	} else if pageSize != 0 {
		qParams.PageSize = pageSize
	}

	// Parse the "sort" query parameter
	if sort := GetQueryParam(r, "sort"); sort != "" {
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
	Movies             []ShowMovieResponse                 `json:"movies"`
	MoviesListMetaData domain.MoviesListPaginationMetadata `json:"pagination_metadata"`
}

func NewListMoviesResponse(movies []*domain.Movie, metadata domain.MoviesListPaginationMetadata) ListMoviesResponse {

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
