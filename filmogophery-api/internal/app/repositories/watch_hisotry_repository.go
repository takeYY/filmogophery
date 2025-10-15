package repositories

import (
	"context"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
)

type (
	IWatchHistoryRepository interface {
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
