package record

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"

	"filmogophery/internal/db"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IQueryRepository interface {
		Find(ctx context.Context) ([]*model.MovieWatchRecord, error)
		FindByImpressionID(ctx context.Context, id *int32) ([]*model.MovieWatchRecord, error)
	}
	ICommandRepository interface {
		Save(ctx context.Context, watchRecord *model.MovieWatchRecord) (*model.MovieWatchRecord, error)
	}

	MovieWatchRecordRepository struct {
		DB *gorm.DB
	}
)

func NewQueryRepository() *IQueryRepository {
	var queryRepo IQueryRepository = &MovieWatchRecordRepository{
		DB: db.READER_DB,
	}
	return &queryRepo
}

func NewCommandRepository() *ICommandRepository {
	var commandRepo ICommandRepository = &MovieWatchRecordRepository{
		DB: db.WRITER_DB,
	}
	return &commandRepo
}

func (r *MovieWatchRecordRepository) Find(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	mwr := query.Use(r.DB).MovieWatchRecord

	movieWatchRecords, err := mwr.WithContext(ctx).Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).Find()
	if err != nil {
		return nil, err
	}

	return movieWatchRecords, nil
}

func (r *MovieWatchRecordRepository) FindByImpressionID(ctx context.Context, id *int32) ([]*model.MovieWatchRecord, error) {
	mwr := query.Use(r.DB).MovieWatchRecord

	movieWatchRecords, err := mwr.WithContext(ctx).
		Preload(
			field.Associations.Scopes(field.RelationFieldUnscoped),
		).
		Where(mwr.MovieImpressionID.Eq(*id)).
		Order(mwr.WatchDate.Desc()).
		Find()
	if err != nil {
		return nil, err
	}

	return movieWatchRecords, nil
}

func (r *MovieWatchRecordRepository) Save(ctx context.Context, watchRecord *model.MovieWatchRecord) (*model.MovieWatchRecord, error) {
	var err error
	defer func() {
		if err != nil {
			r.DB.Rollback()
		} else {
			r.DB.Commit()
		}
	}()

	mwr := query.Use(r.DB).MovieWatchRecord

	err = mwr.WithContext(ctx).Create(watchRecord)

	return watchRecord, err
}
