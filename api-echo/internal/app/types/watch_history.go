package types

import (
	"time"

	"filmogophery/internal/pkg/constant"
)

type (
	MovieWatchHistory struct {
		ID        int32     `json:"id"`
		Platform  Platform  `json:"platform"`
		WatchedAt time.Time `json:"watchedAt"`
	}

	WatchHistory struct {
		ID        int32         `json:"id"`
		WatchedAt constant.Date `json:"watchedAt"`
		Platform  Platform      `json:"platform"`
		Movie     Movie         `json:"movie"`
	}
)
