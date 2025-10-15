package repositories

import (
	"context"

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
