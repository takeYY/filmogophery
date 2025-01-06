package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IMediaRepository interface {
		// --- Read --- //

		// 鑑賞媒体を全て取得
		FindAll(ctx context.Context) ([]*model.WatchMedia, error)
		// コードと一致するメディアを取得
		FindByCode(ctx context.Context, code *string) (*model.WatchMedia, error)
	}

	mediaRepository struct {
		ReaderDB *gorm.DB
		// WriterDB *gorm.DB
	}
)

func NewMediaRepository(db *gorm.DB) *IMediaRepository {
	var repo IMediaRepository = &mediaRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		// WriterDB: db.Clauses(dbresolver.Write),
	}
	return &repo
}

func (r *mediaRepository) FindAll(ctx context.Context) ([]*model.WatchMedia, error) {
	wm := query.Use(r.ReaderDB).WatchMedia

	watchMedia, err := wm.WithContext(ctx).
		Order(wm.ID).
		Find()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return make([]*model.WatchMedia, 0), nil
	}

	return watchMedia, err
}

func (r *mediaRepository) FindByCode(ctx context.Context, code *string) (*model.WatchMedia, error) {
	wm := query.Use(r.ReaderDB).WatchMedia

	watchMedia, err := wm.WithContext(ctx).
		Where(wm.Code.Eq(*code)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return nil, nil
	}

	return watchMedia, err
}
