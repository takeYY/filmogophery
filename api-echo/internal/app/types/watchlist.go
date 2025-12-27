package types

import "time"

type (
	Watchlist struct {
		ID       int32     `json:"id"`
		AddedAt  time.Time `json:"added_at"`
		Priority int32     `json:"priority"`
		Movie    Movie     `json:"movie"`
	}
)
