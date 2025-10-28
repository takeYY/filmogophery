package review

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/constant"
)

type (
	CreateReviewHistoryUseCase interface {
		Run(ctx context.Context, reviewID int32, platformID int32, watchedDate *constant.Date) error
	}

	createReviewHistoryInteractor struct {
		reviewService   services.IReviewService
		platformService services.IPlatformService
	}
)

func NewCreateReviewHistoryInteractor(
	reviewService services.IReviewService,
	platformService services.IPlatformService,
) CreateReviewHistoryUseCase {
	return &createReviewHistoryInteractor{
		reviewService,
		platformService,
	}
}

func (i *createReviewHistoryInteractor) Run(
	ctx context.Context, reviewID int32, platformID int32, watchedDate *constant.Date,
) error {
	// レビューの存在確認
	review, err := i.reviewService.GetReviewByID(ctx, 1, reviewID)
	if err != nil {
		return err
	}

	// プラットフォームの存在確認
	platform, err := i.platformService.GetByID(ctx, platformID)
	if err != nil {
		return err
	}

	// 視聴履歴登録
	return i.reviewService.CreateWatchHistory(ctx, nil, review, platform, watchedDate)
}
