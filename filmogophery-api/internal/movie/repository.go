package movie

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		FindByID(ctx context.Context, id *int32) (*model.Movie, error)
		Find(ctx context.Context) ([]*model.Movie, error)
	}
	ICommandRepository interface {
		Save(movie *model.Movie) (*model.Movie, error)
	}

	MovieRepository struct {
		DB *gorm.DB
	}
)

func (r *MovieRepository) FindByID(ctx context.Context, id *int32) (*model.Movie, error) {
	m := query.Use(r.DB).Movie

	movie, err := m.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
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
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Find()
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (r *MovieRepository) Save(movie *model.Movie) (*model.Movie, error) {
	var err error
	defer func() {
		if err != nil {
			r.DB.Rollback()
		} else {
			r.DB.Commit()
		}
	}()

	m := query.Use(r.DB).Movie
	err = m.Create(movie)
	if err != nil {
		return nil, err
	}

	return movie, nil
}
