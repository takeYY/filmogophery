package movie

type (
	CreateMovieDto struct {
		Title       string   `json:"title"`
		Overview    *string  `json:"overview"`
		ReleaseDate string   `json:"release_date"`
		RunTime     int32    `json:"run_time"`
		TmdbID      *int32   `json:"tmdb_id"`
		Genres      []string `json:"genres"`
	}
)
