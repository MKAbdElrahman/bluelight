package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"bluelight.mkcodedev.com/src/api/handlers/jsonio"
	"bluelight.mkcodedev.com/src/core/domain"
)

func newShowMovieHandlerFunc(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parsedId, err := parseIdFromPath(r)
		if err != nil || parsedId < 1 {
			notFound(logger, w, r)
			return
		}

		m := domain.Movie{
			Id:               int64(parsedId),
			CreatedAt:        time.Now(),
			Title:            "Casablanca",
			RuntimeInMinutes: 102,
			Genres:           []string{"drama", "romance", "war"},
			Version:          1,
		}

		err = jsonio.SendJSON(w, jsonio.Envelope{"movie": m}, http.StatusOK, nil)
		if err != nil {
			internal(logger, w, r, err)
			return
		}

	}
}

func parseIdFromPath(r *http.Request) (int, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.Atoi(idFromPath)
	return parsedId, err
}
