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

	CreateMovieRecordDto struct {
		MovieID      int32   `json:"movieId"`
		ImpressionID int32   `json:"impressionId"`
		Media        string  `json:"media"`
		WatchDate    string  `json:"watchDate"`
		Rating       float32 `json:"rating"`
		Note         string  `json:"note"`
	}
)
