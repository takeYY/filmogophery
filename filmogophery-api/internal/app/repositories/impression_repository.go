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
	IImpressionRepository interface {
		// --- Create --- //

		// 感想を作成
		Save(ctx context.Context, input SaveImpressionInput) error

		// --- Read --- //

		// 感想を全て取得
		FindAll(ctx context.Context) ([]*model.MovieImpression, error)

		// --- Update --- //

		// 感想を更新
		Update(ctx context.Context, input UpdateImpressionInput) (int64, error)
	}
	SaveImpressionInput struct {
		Target *model.MovieImpression
		Tx     *gorm.DB
	}
	UpdateImpressionInput struct {
		Target *model.MovieImpression
		Tx     *gorm.DB
	}

	impressionRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewImpressionRepository(db *gorm.DB) *IImpressionRepository {
	var repo IImpressionRepository = &impressionRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
	return &repo
}

func (r *impressionRepository) Save(ctx context.Context, input SaveImpressionInput) error {
	mi := query.Use(r.WriterDB).MovieImpression
	if input.Tx != nil {
		mi = query.Use(input.Tx).MovieImpression
	}

	return mi.Create(input.Target)
}

func (r *impressionRepository) FindAll(ctx context.Context) ([]*model.MovieImpression, error) {
	mi := query.Use(r.ReaderDB).MovieImpression

	movieImpressions, err := mi.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Find()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return make([]*model.MovieImpression, 0), nil
	}

	return movieImpressions, err
}

func (r *impressionRepository) Update(ctx context.Context, input UpdateImpressionInput) (int64, error) {
	mi := query.Use(r.WriterDB).MovieImpression
	if input.Tx != nil {
		mi = query.Use(input.Tx).MovieImpression
	}

	result, err := mi.WithContext(ctx).
		Where(mi.ID.Eq(input.Target.ID)).
		Updates(input.Target)

	return result.RowsAffected, err
}
