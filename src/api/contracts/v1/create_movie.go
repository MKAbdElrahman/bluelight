package v1

type CreateMovieRequest struct {
	Body struct {
		Title   string   `json:"title"`
		Year    int      `json:"year"`
		Runtime string   `json:"runtime"`
		Genres  []string `json:"genres"`
	}
}
