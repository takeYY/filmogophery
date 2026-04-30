package review

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/gen/model"
)

type (
	CreateReviewUseCase interface {
		Run(ctx context.Context, operator *model.Users, movieID int32, rating *float64, comment *string) error
	}

	createReviewInteractor struct {
		movieService  services.IMovieService
		reviewService services.IReviewService
		pointService  services.IPointService
	}
)

func NewCreateReviewInteractor(
	movieService services.IMovieService,
	reviewService services.IReviewService,
	pointService services.IPointService,
) CreateReviewUseCase {
	return &createReviewInteractor{
		movieService,
		reviewService,
		pointService,
	}
}

func (i *createReviewInteractor) Run(ctx context.Context, operator *model.Users, movieID int32, rating *float64, comment *string) error {
	// 映画の存在確認
	movie, err := i.movieService.GetMovieByID(ctx, movieID)
	if err != nil {
		return err
	}

	// レビューを作成
	if err := i.reviewService.CreateReview(ctx, nil, operator, movie, rating, comment); err != nil {
		return err
	}

	// 作成したレビューを取得してポイント付与
	review, err := i.reviewService.GetReviewByMovieID(ctx, operator, movie)
	if err != nil {
		return err
	}

	return i.pointService.GrantReviewPoints(ctx, nil, operator, review)
}
