package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IPlatformRepository interface {
		// --- Read --- //

		// 全てのプラットフォームを取得
		FindAll(ctx context.Context) ([]*model.Platforms, error)
		// IDに一致するプラットフォームを取得
		FindByID(ctx context.Context, id int32) (*model.Platforms, error)
	}

	platformRepository struct {
		ReaderDB *gorm.DB
		// WriterDB *gorm.DB
	}
)

func NewPlatformRepository(db *gorm.DB) IPlatformRepository {
	return &platformRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		// NOTE: 書き込みは不要
	}
}

// 全てのプラットフォームを取得
func (r *platformRepository) FindAll(ctx context.Context) ([]*model.Platforms, error) {
	p := query.Use(r.ReaderDB).Platforms

	return p.WithContext(ctx).Find()
}

// IDに一致するプラットフォームを取得
func (r *platformRepository) FindByID(ctx context.Context, id int32) (*model.Platforms, error) {
	p := query.Use(r.ReaderDB).Platforms

	result, err := p.WithContext(ctx).
		Where(p.ID.Eq(id)).
		Take()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}

	return result, nil
}
