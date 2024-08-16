package media

import (
	"context"

	"gorm.io/gorm"

	"filmogophery/internal/db"
	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		Find(ctx context.Context) ([]*model.WatchMedia, error)
		GetMediaIdByCode(ctx context.Context, code *string) (*int32, error)
	}

	WatchMediaRepository struct {
		DB *gorm.DB
	}
)

func NewQueryRepository() *IQueryRepository {
	var queryRepo IQueryRepository = &WatchMediaRepository{
		DB: db.READER_DB,
	}
	return &queryRepo
}

func (r *WatchMediaRepository) Find(ctx context.Context) ([]*model.WatchMedia, error) {
	wm := query.Use(r.DB).WatchMedia

	watchMedia, err := wm.WithContext(ctx).Order(wm.ID).Find()
	if err != nil {
		return nil, err
	}

	return watchMedia, nil
}

func (r *WatchMediaRepository) GetMediaIdByCode(ctx context.Context, code *string) (*int32, error) {
	wm := query.Use(r.DB).WatchMedia

	watchMedia, err := wm.WithContext(ctx).Where(wm.Code.Eq(*code)).First()
	if err != nil {
		return nil, err
	}
	return &watchMedia.ID, nil
}
