package impression

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		Find(ctx context.Context) ([]*model.MovieImpression, error)
	}

	MovieImpressionRepository struct {
		DB *gorm.DB
	}
)

func (r *MovieImpressionRepository) Find(ctx context.Context) ([]*model.MovieImpression, error) {
	mi := query.Use(r.DB).MovieImpression

	movieImpressions, err := mi.WithContext(ctx).Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).Find()
	if err != nil {
		return nil, err
	}

	return movieImpressions, nil
}
