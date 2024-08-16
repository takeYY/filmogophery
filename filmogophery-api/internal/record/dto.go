package record

type (
	CreateMovieRecordDto struct {
		MovieID      int32   `json:"movieId"`
		ImpressionID int32   `json:"impressionId"`
		Media        string  `json:"media"`
		WatchDate    string  `json:"watchDate"`
		Rating       float32 `json:"rating"`
		Note         string  `json:"note"`
	}
)
