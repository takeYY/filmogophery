package repositories

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IGenreRepository interface {
		// --- Read --- //

		// 全てのジャンルを取得
		FindAll(ctx context.Context) ([]*model.Genres, error)
		// 名前と一致するジャンルを取得
		FindByNames(ctx context.Context, names []string) ([]*model.Genres, error)
	}

	genreRepository struct {
		ReaderDB *gorm.DB
		// WriterDB *gorm.DB
	}
)

func NewGenreRepository(db *gorm.DB) IGenreRepository {
	return &genreRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		// NOTE: 書き込みは不要
	}
}

// 全てのジャンルを取得
func (r *genreRepository) FindAll(ctx context.Context) ([]*model.Genres, error) {
	g := query.Use(r.ReaderDB).Genres

	return g.WithContext(ctx).Find()
}

func (r *genreRepository) FindByNames(ctx context.Context, names []string) ([]*model.Genres, error) {
	g := query.Use(r.ReaderDB).Genres

	return g.WithContext(ctx).
		Where(g.Name.In(names...)).
		Find()
}
