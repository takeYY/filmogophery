package review

import (
	"context"

	"gorm.io/gorm"

	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
)

// WatchHistoryInput はレビュー登録時に同時に登録する視聴履歴の入力
type WatchHistoryInput struct {
	PlatformID  int32
	WatchedDate *constant.Date
}

type (
	CreateReviewUseCase interface {
		Run(ctx context.Context, operator *model.Users, movieID int32, rating *float64, comment *string, watchHistory *WatchHistoryInput) error
	}

	createReviewInteractor struct {
		db              *gorm.DB
		movieService    services.IMovieService
		reviewService   services.IReviewService
		platformService services.IPlatformService
		pointService    services.IPointService
	}
)

func NewCreateReviewInteractor(
	db *gorm.DB,
	movieService services.IMovieService,
	reviewService services.IReviewService,
	platformService services.IPlatformService,
	pointService services.IPointService,
) CreateReviewUseCase {
	return &createReviewInteractor{
		db,
		movieService,
		reviewService,
		platformService,
		pointService,
	}
}

func (i *createReviewInteractor) Run(
	ctx context.Context, operator *model.Users, movieID int32, rating *float64, comment *string, watchHistory *WatchHistoryInput,
) error {
	// 映画の存在確認
	movie, err := i.movieService.GetMovieByID(ctx, movieID)
	if err != nil {
		return err
	}

	// プラットフォームの存在確認（視聴履歴あり時のみ）
	var platform *model.Platforms
	if watchHistory != nil {
		platform, err = i.platformService.GetByID(ctx, watchHistory.PlatformID)
		if err != nil {
			return err
		}
	}

	// トランザクション内でレビュー・視聴履歴・ポイントをまとめて登録
	return i.db.Transaction(func(tx *gorm.DB) error {
		// レビューを作成（INSERT後にIDが書き戻される）
		review, err := i.reviewService.CreateReview(ctx, tx, operator, movie, rating, comment)
		if err != nil {
			return err
		}

		// レビューポイントを付与
		if err := i.pointService.GrantReviewPoints(ctx, tx, operator, review); err != nil {
			return err
		}

		// 視聴履歴を登録（入力がある場合のみ）
		if watchHistory != nil {
			wh, err := i.reviewService.CreateWatchHistory(ctx, tx, operator, review, platform, watchHistory.WatchedDate)
			if err != nil {
				return err
			}

			// 視聴履歴ポイントを付与
			if err := i.pointService.GrantWatchHistoryPoints(ctx, tx, operator, wh, movie); err != nil {
				return err
			}
		}

		return nil
	})
}
