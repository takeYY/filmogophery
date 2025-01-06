package types

type (
	Movie struct {
		ID          int32   `json:"id"`
		Title       string  `json:"title"`
		Overview    string  `json:"overview"`
		ReleaseDate string  `json:"releaseDate"`
		RunTime     int32   `json:"runTime"`
		PosterURL   *string `json:"posterURL"`
		TmdbID      int32   `json:"tmdbID"`
		Genres      []Genre `json:"genres"`
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
