package review

import (
	"context"

	"filmogophery/internal/app/services"
)

type (
	CreateReviewUseCase interface {
		Run(ctx context.Context, movieID int32, rating *float64, comment *string) error
	}

	createReviewInteractor struct {
		movieService  services.IMovieService
		reviewService services.IReviewService
	}
)

func NewCreateReviewInteractor(
	movieService services.IMovieService,
	reviewService services.IReviewService,
) CreateReviewUseCase {
	return &createReviewInteractor{
		movieService,
		reviewService,
	}
}

func (i *createReviewInteractor) Run(ctx context.Context, movieID int32, rating *float64, comment *string) error {
	// 映画の存在確認
	movie, err := i.movieService.GetMovieByID(ctx, movieID)
	if err != nil {
		return err
	}

	// レビューを作成
	err = i.reviewService.CreateReview(ctx, nil, 1, movie, rating, comment)
	if err != nil {
		return err
	}

	return nil
}
