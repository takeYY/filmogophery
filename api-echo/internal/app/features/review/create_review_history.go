package review

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
)

type (
	CreateReviewHistoryUseCase interface {
		Run(ctx context.Context, operator *model.Users, reviewID int32, platformID int32, watchedDate *constant.Date) error
	}

	createReviewHistoryInteractor struct {
		reviewService   services.IReviewService
		platformService services.IPlatformService
		pointService    services.IPointService
		movieService    services.IMovieService
	}
)

func NewCreateReviewHistoryInteractor(
	reviewService services.IReviewService,
	platformService services.IPlatformService,
	pointService services.IPointService,
	movieService services.IMovieService,
) CreateReviewHistoryUseCase {
	return &createReviewHistoryInteractor{
		reviewService,
		platformService,
		pointService,
		movieService,
	}
}

func (i *createReviewHistoryInteractor) Run(
	ctx context.Context, operator *model.Users, reviewID int32, platformID int32, watchedDate *constant.Date,
) error {
	// レビューの存在確認
	review, err := i.reviewService.GetReviewByID(ctx, operator, reviewID)
	if err != nil {
		return err
	}

	// プラットフォームの存在確認
	platform, err := i.platformService.GetByID(ctx, platformID)
	if err != nil {
		return err
	}

	// 視聴履歴登録
	watchHistory, err := i.reviewService.CreateWatchHistory(ctx, nil, operator, review, platform, watchedDate)
	if err != nil {
		return err
	}

	// 映画情報を取得してポイント付与
	movie, err := i.movieService.GetMovieByID(ctx, review.MovieID)
	if err != nil {
		return err
	}

	return i.pointService.GrantWatchHistoryPoints(ctx, nil, operator, watchHistory, movie)
}
