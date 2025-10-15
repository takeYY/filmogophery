package types

import (
	"time"

	"filmogophery/internal/pkg/gen/model"
)

type (
	Review struct {
		ID        int32     `json:"id"`
		Rating    *float64  `json:"rating"`
		Comment   *string   `json:"comment"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	ReviewHistory struct {
		ID        int32     `json:"id"`
		Platform  Platform  `json:"platform"`
		WatchedAt time.Time `json:"watchedAt"`
	}
)

func NewReviewByModel(review *model.Reviews) *Review {
	if review == nil {
		return nil
	}

	return &Review{
		ID:        review.ID,
		Rating:    review.Rating,
		Comment:   review.Comment,
		CreatedAt: *review.CreatedAt,
		UpdatedAt: *review.UpdatedAt,
	}
}
