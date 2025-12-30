package types

import (
	"filmogophery/internal/pkg/constant"
)

type (
	Movie struct {
		ID             int32         `json:"id"`
		Title          string        `json:"title"`
		Overview       string        `json:"overview"`
		ReleaseDate    constant.Date `json:"releaseDate"`
		RuntimeMinutes int32         `json:"runtimeMinutes"`
		PosterURL      *string       `json:"posterURL"`
		TmdbID         int32         `json:"tmdbID"`
		Genres         []Genre       `json:"genres"`
	}
	MovieDetail struct {
		VoteAverage float32 `json:"voteAverage"`
		VoteCount   int32   `json:"voteCount"`
		Series      *Series `json:"series"`
		Review      *Review `json:"review"`
		Movie
	}

	SearchMovie struct {
		TmdbID      int32    `json:"tmdbID"`
		Title       string   `json:"title"`
		Overview    string   `json:"overview"`
		Popularity  int32    `json:"popularity"`
		PosterURL   string   `json:"posterURL"`
		ReleaseDate string   `json:"releaseDate"`
		VoteAverage int32    `json:"voteAverage"`
		VoteCount   int32    `json:"voteCount"`
		Genres      []string `json:"genres"`
	}
)
