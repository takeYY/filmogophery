package record

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"filmogophery/pkg/gen/model"
	"filmogophery/pkg/gen/query"
)

type (
	IQueryRepository interface {
		Find(ctx context.Context) ([]*model.MovieWatchRecord, error)
	}

	MovieWatchRecordRepository struct {
		DB *gorm.DB
	}
)

func (r *MovieWatchRecordRepository) Find(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	mwr := query.Use(r.DB).MovieWatchRecord

	movieWatchRecords, err := mwr.WithContext(ctx).Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).Find()
	if err != nil {
		return nil, err
	}

	return movieWatchRecords, nil
}
