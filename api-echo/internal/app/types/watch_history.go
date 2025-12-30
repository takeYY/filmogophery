package types

import "time"

type (
	MovieWatchHistory struct {
		ID        int32     `json:"id"`
		Platform  Platform  `json:"platform"`
		WatchedAt time.Time `json:"watchedAt"`
	}
)
