package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"time"

	v1 "bluelight.mkcodedev.com/src/api/contracts/v1"
	"bluelight.mkcodedev.com/src/core/domain"
	"bluelight.mkcodedev.com/src/lib/jsonio"
)

func newShowMovieHandlerFunc(logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parsedId, err := parseIdFromPath(r)
		if err != nil || parsedId < 1 {
			sendClientError(logger, w, r, v1.NotFoundError)
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
			sendServerError(logger, w, r, v1.InternalServerError)
			return
		}

	}
}

func parseIdFromPath(r *http.Request) (int, error) {
	idFromPath := r.PathValue("id")
	parsedId, err := strconv.Atoi(idFromPath)
	return parsedId, err
}
