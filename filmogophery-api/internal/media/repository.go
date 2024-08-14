package media

import (
	"context"

	"gorm.io/gorm"

	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		Find(ctx context.Context) ([]*model.WatchMedia, error)
	}

	WatchMediaRepository struct {
		DB *gorm.DB
	}
)

func (r *WatchMediaRepository) Find(ctx context.Context) ([]*model.WatchMedia, error) {
	wm := query.Use(r.DB).WatchMedia

	watchMedia, err := wm.WithContext(ctx).Order(wm.ID).Find()
	if err != nil {
		return nil, err
	}

	return watchMedia, nil
}
