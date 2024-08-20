package domain

import "time"

type Movie struct {
	Id               int64     `json:"id"`
	CreatedAt        time.Time `json:"-"`
	Title            string    `json:"title"`
	Year             int32     `json:"year,omitempty"`
	RuntimeInMinutes int32     `json:"runtime,omitempty"`
	Genres           []string  `json:"genres,omitempty"`
	Version          int32     `json:"version"`
}
