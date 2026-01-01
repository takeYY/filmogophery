package review

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetReviewHistoryUseCase interface {
		Run(ctx context.Context, operator *model.Users, reviewID int32) ([]*types.ReviewHistory, error)
	}

	getReviewHistoryInteractor struct {
		reviewService services.IReviewService
	}
)

func NewGetReviewHistoryInteractor(
	reviewService services.IReviewService,
) GetReviewHistoryUseCase {
	return &getReviewHistoryInteractor{
		reviewService,
	}
}

func (i *getReviewHistoryInteractor) Run(ctx context.Context, operator *model.Users, reviewID int32) ([]*types.ReviewHistory, error) {
	// レビューの存在確認
	review, err := i.reviewService.GetReviewByID(ctx, 1, reviewID)
	if err != nil {
		return nil, err
	}

	// 視聴履歴を取得
	watchHistories, err := i.reviewService.GetWatchHistoryByReviewID(ctx, operator, review)
	if err != nil {
		return nil, err
	}

	response := make([]*types.ReviewHistory, 0, len(watchHistories))
	for _, wh := range watchHistories {
		response = append(response, &types.ReviewHistory{
			ID:        wh.ID,
			Platform:  types.NewPlatformByModel(wh.Platform),
			WatchedAt: *wh.WatchedDate,
		})
	}

	return response, nil
}
