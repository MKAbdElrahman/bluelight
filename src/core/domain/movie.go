package domain

import "time"

type Movie struct {
	Id               int64
	CreatedAt        time.Time
	Title            string
	Year             int32
	RuntimeInMinutes int32
	Genres           []string
	Version          int32
}
