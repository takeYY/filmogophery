package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IWatchHistoryRepository interface {
		// --- Create --- //

		// 視聴履歴を作成
		Save(ctx context.Context, tx *gorm.DB, watchHistory *model.WatchHistory) error

		// --- Read --- //

		// レビューIDに一致する視聴履歴を取得
		FindByReviewID(ctx context.Context, reviewID int32) ([]*model.WatchHistory, error)
	}
	watchHistoryRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewWatchHistoryRepository(db *gorm.DB) IWatchHistoryRepository {
	return &watchHistoryRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// 視聴履歴を作成
func (r *watchHistoryRepository) Save(ctx context.Context, tx *gorm.DB, watchHistory *model.WatchHistory) error {
	wh := query.Use(r.ReaderDB).WatchHistory
	if tx != nil {
		wh = query.Use(tx).WatchHistory
	}

	return wh.WithContext(ctx).Create(watchHistory)
}

// レビューIDに一致する視聴履歴を取得
func (r *watchHistoryRepository) FindByReviewID(
	ctx context.Context, reviewID int32,
) ([]*model.WatchHistory, error) {
	wh := query.Use(r.ReaderDB).WatchHistory

	return wh.WithContext(ctx).
		Preload(wh.Platform).
		Where(wh.ReviewID.Eq(reviewID)).
		Find()
}
