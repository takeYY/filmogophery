package types

type (
	Impression struct {
		ID      int32    `json:"id"`
		Status  string   `json:"status"`
		Rating  *float32 `json:"rating"`
		Note    *string  `json:"note"`
		Records []Record `json:"records"`
	}

	// NOTE: 今後要らなくなるので頃合いを見て消すこと
	MovieImpression struct {
		ID           int32  `json:"id"`
		MovieID      int32  `json:"movieID"`
		Status       bool   `json:"status"`
		Rating       int32  `json:"rating"`
		Note         string `json:"note"`
		Movie        *Movie `json:"movie"`
		WatchRecords *int32 `json:"watch_records"`
	}
)
