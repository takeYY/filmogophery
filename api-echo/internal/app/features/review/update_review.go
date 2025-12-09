package review

import (
	"context"

	"filmogophery/internal/app/services"
)

type (
	UpdateReviewUseCase interface {
		Run(ctx context.Context, reviewID int32, rating *float64, comment *string) error
	}

	updateReviewInteractor struct {
		reviewService services.IReviewService
	}
)

func NewUpdateReviewInteractor(
	reviewService services.IReviewService,
) UpdateReviewUseCase {
	return &updateReviewInteractor{
		reviewService,
	}
}

func (i *updateReviewInteractor) Run(ctx context.Context, reviewID int32, rating *float64, comment *string) error {
	// レビューの存在確認
	review, err := i.reviewService.GetReviewByID(ctx, 1, reviewID)
	if err != nil {
		return err
	}

	// レビューを更新
	err = i.reviewService.UpdateReview(ctx, nil, review, rating, comment)
	if err != nil {
		return err
	}

	return nil
}
