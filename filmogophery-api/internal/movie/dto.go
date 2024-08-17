package movie

type (
	CreateMovieDto struct {
		TmdbID int32 `json:"tmdbID"`
		Status bool  `json:"status"`
	}

	SeriesDetailDto struct {
		Name      string `json:"name"`
		PosterURL string `json:"posterURL"`
	}

	ImpressionDetailDto struct {
		ID     int32    `json:"id"`
		Status bool     `json:"status"`
		Rating *float32 `json:"rating"`
		Note   *string  `json:"note"`
	}

	WatchRecordDetailDto struct {
		WatchDate  string `json:"watchDate"`
		WatchMedia string `json:"watchMedia"`
	}

	MovieDetailDto struct {
		ID           int32                   `json:"id"`
		Title        string                  `json:"title"`
		Overview     *string                 `json:"overview"`
		ReleaseDate  string                  `json:"releaseDate"`
		RunTime      int32                   `json:"runTime"`
		Genres       []*string               `json:"genres"`
		PosterURL    string                  `json:"posterURL"`
		VoteAverage  float32                 `json:"voteAverage"`
		VoteCount    int32                   `json:"voteCount"`
		Series       *SeriesDetailDto        `json:"series"`
		Impression   *ImpressionDetailDto    `json:"impression"`
		WatchRecords []*WatchRecordDetailDto `json:"watchRecords"`
	}
)
