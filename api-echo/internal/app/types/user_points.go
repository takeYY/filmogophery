package types

import (
	"time"

	"filmogophery/internal/pkg/constant"
)

type (
	PointHistoryItem struct {
		ID          int32                `json:"id"`
		Points      int32                `json:"points"`
		Action      constant.PointAction `json:"action"`
		ReferenceID int32                `json:"referenceId"`
		CreatedAt   *time.Time           `json:"createdAt"`
	}

	UserPoints struct {
		TotalPoints       int32              `json:"totalPoints"`
		Level             int32              `json:"level"`
		NextLevelPoints   int32              `json:"nextLevelPoints"`
		CurrentLevelWidth int32              `json:"currentLevelWidth"`
		PointHistory      []PointHistoryItem `json:"pointHistory"`
	}
)
