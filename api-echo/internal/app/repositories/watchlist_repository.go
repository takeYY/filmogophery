package repositories

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IWatchlistRepository interface {
		// --- Create --- //

		// ウォッチリストを登録
		Create(ctx context.Context, tx *gorm.DB, watchlist *model.Watchlist) error

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

// ウォッチリストを登録
func (r *watchlistRepository) Create(ctx context.Context, tx *gorm.DB, watchlist *model.Watchlist) error {
	wl := query.Use(r.WriterDB).Watchlist
	if tx != nil {
		wl = query.Use(tx).Watchlist
	}

	return wl.WithContext(ctx).
		Omit(field.AssociationFields).
		Create(watchlist)
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
