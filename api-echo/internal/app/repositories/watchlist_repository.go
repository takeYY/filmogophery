package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IWatchlistRepository interface {
		// --- Create --- //

		// --- Read --- //

		// ユーザーのウォッチリストを取得
		FindByUserID(ctx context.Context, user *model.Users, limit int32, offset int32) ([]*model.Watchlist, error)

		// --- Update --- //

		// --- Delete --- //
	}

	watchlistRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewWatchlistRepository(db *gorm.DB) IWatchlistRepository {
	return &watchlistRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// ユーザーのウォッチリストを取得
func (r *watchlistRepository) FindByUserID(
	ctx context.Context, user *model.Users, limit int32, offset int32,
) ([]*model.Watchlist, error) {
	wl := query.Use(r.ReaderDB).Watchlist

	return wl.WithContext(ctx).
		Preload(wl.Movie.Genres).
		Where(wl.UserID.Eq(user.ID)).
		Order(wl.AddedAt.Desc()).
		Limit(int(limit)).
		Offset(int(offset)).
		Find()
}
