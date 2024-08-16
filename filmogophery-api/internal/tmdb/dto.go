package tmdb

type (
	SearchMovieDto struct {
		TmdbID      int32     `json:"tmdbID"`
		Title       string    `json:"title"`
		Overview    string    `json:"overview"`
		Popularity  float32   `json:"popularity"`
		PosterURL   string    `json:"posterURL"`
		ReleaseDate string    `json:"releaseDate"`
		VoteAverage float32   `json:"voteAverage"`
		VoteCount   int32     `json:"voteCount"`
		Genres      []*string `json:"genres"`
	}
)
