package movie

type (
	GetMovieDetailRequest struct {
		ID int32 `param:"id"`
	}

	PostMovieImpression struct {
		ID        int32   `param:"id"`
		WatchDate string  `json:"watchDate"`
		MediaCode string  `json:"mediaCode"`
		Rating    float32 `json:"rating"`
		Note      string  `json:"note"`
	}

	PutMovieImpression struct {
		ID     int32   `param:"id"`
		Rating float32 `json:"rating"`
		Note   string  `json:"note"`
	}

	PutMovieRecord struct {
		ID        int32  `param:"id"`
		RecordID  int32  `param:"recordId"`
		Date      string `json:"date"`
		MediaCode string `json:"mediaCode"`
	}
)
