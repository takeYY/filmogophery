package movie

import (
	"context"

	"gorm.io/gorm"

	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		FindByID(ctx context.Context, id *int64) (*model.Movie, error)
		Find(ctx context.Context) ([]*model.Movie, error)
	}

	MovieRepository struct {
		DB *gorm.DB
	}
)

func (r *MovieRepository) FindByID(ctx context.Context, id *int64) (*model.Movie, error) {
	m := query.Use(r.DB).Movie

	movie, err := m.WithContext(ctx).
		Preload(m.Genres).
		Preload(m.Poster).
		Preload(m.Series).
		Where(m.ID.Eq(*id)).
		First()
	if err != nil {
		return nil, err
	}

	return movie, nil
}

func (r *MovieRepository) Find(ctx context.Context) ([]*model.Movie, error) {
	m := query.Use(r.DB).Movie

	movies, err := m.WithContext(ctx).
		Preload(m.Genres).
		Preload(m.Poster).
		Preload(m.Series).
		Find()
	if err != nil {
		return nil, err
	}

	return movies, nil
}
