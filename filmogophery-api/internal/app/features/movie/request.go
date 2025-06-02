package movie

type (
	GetMovieDetailRequest struct {
		ID int32 `param:"id"`
	}

	PostMovieImpression struct {
		ID         int32   `param:"id"`
		WatchDate  string  `json:"watchDate"`
		WatchMedia string  `json:"watchMedia"`
		Rating     float32 `json:"rating"`
		Note       string  `json:"note"`
	}
)
