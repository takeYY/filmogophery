package types

type (
	TrendingMovie struct {
		ID        int32   `json:"id"`
		Title     string  `json:"title"`
		PosterURL *string `json:"posterURL"`
		TmdbID    int32   `json:"tmdbID"`
	}
)
