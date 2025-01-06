package repositories

import (
	"context"
	"errors"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IRecordRepository interface {
		// --- Create --- //

		// 鑑賞記録を作成
		Save(ctx context.Context, input SaveRecordInput) error

		// --- Read --- //

		// 全ての鑑賞記録を取得
		FindAll(ctx context.Context) ([]*model.MovieWatchRecord, error)
		// 感想IDと一致する鑑賞記録を取得
		FindByImpressionID(ctx context.Context, id *int32) ([]*model.MovieWatchRecord, error)
	}
	SaveRecordInput struct {
		Target *model.MovieWatchRecord
		Tx     *gorm.DB
	}

	recordRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewRecordRepository(db *gorm.DB) *IRecordRepository {
	var repo IRecordRepository = &recordRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
	return &repo
}

func (r *recordRepository) Save(ctx context.Context, input SaveRecordInput) error {
	mwr := query.Use(r.WriterDB).MovieWatchRecord
	if input.Tx != nil {
		mwr = query.Use(input.Tx).MovieWatchRecord
	}

	return mwr.WithContext(ctx).Create(input.Target)
}

func (r *recordRepository) FindAll(ctx context.Context) ([]*model.MovieWatchRecord, error) {
	mwr := query.Use(r.ReaderDB).MovieWatchRecord

	movieWatchRecords, err := mwr.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Find()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return make([]*model.MovieWatchRecord, 0), nil
	}

	return movieWatchRecords, err
}

func (r *recordRepository) FindByImpressionID(ctx context.Context, id *int32) ([]*model.MovieWatchRecord, error) {
	mwr := query.Use(r.ReaderDB).MovieWatchRecord

	movieWatchRecords, err := mwr.WithContext(ctx).
		Preload(
			field.Associations.Scopes(field.RelationFieldUnscoped),
		).
		Where(mwr.MovieImpressionID.Eq(*id)).
		Order(mwr.WatchDate.Desc()).
		Find()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return make([]*model.MovieWatchRecord, 0), nil
	}

	return movieWatchRecords, err
}
