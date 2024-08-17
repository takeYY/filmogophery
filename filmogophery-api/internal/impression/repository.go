package impression

import (
	"context"

	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"

	"filmogophery/internal/db"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IQueryRepository interface {
		Find(ctx context.Context) ([]*model.MovieImpression, error)
	}
	ICommandRepository interface {
		UpdateImpression(ctx context.Context, movieImpression *model.MovieImpression) (*gen.ResultInfo, error)
	}

	MovieImpressionRepository struct {
		DB *gorm.DB
	}
)

func NewQueryRepository() *IQueryRepository {
	var queryRepo IQueryRepository = &MovieImpressionRepository{
		DB: db.READER_DB,
	}
	return &queryRepo
}

func NewCommandRepository() *ICommandRepository {
	var commandRepo ICommandRepository = &MovieImpressionRepository{
		DB: db.WRITER_DB,
	}
	return &commandRepo
}

func (r *MovieImpressionRepository) Find(ctx context.Context) ([]*model.MovieImpression, error) {
	mi := query.Use(r.DB).MovieImpression

	movieImpressions, err := mi.WithContext(ctx).Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).Find()
	if err != nil {
		return nil, err
	}

	return movieImpressions, nil
}

func (r *MovieImpressionRepository) UpdateImpression(ctx context.Context, movieImpression *model.MovieImpression) (*gen.ResultInfo, error) {
	var err error
	defer func() {
		if err != nil {
			r.DB.Rollback()
		} else {
			r.DB.Commit()
		}
	}()

	mi := query.Use(r.DB).MovieImpression

	var result gen.ResultInfo
	result, err = mi.WithContext(ctx).Where(mi.ID.Eq(movieImpression.ID)).Updates(movieImpression)

	return &result, err
}
