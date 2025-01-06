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
	IGenreRepository interface {
		// --- Read --- //

		// 名前と一致するジャンルを取得
		FindByNames(ctx context.Context, names []string) ([]*model.Genre, error)
	}

	genreRepository struct {
		ReaderDB *gorm.DB
		// WriterDB *gorm.DB
	}
)

func NewGenreRepository(db *gorm.DB) *IGenreRepository {
	var repo IGenreRepository = &genreRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		// WriterDB: db.Clauses(dbresolver.Write),
	}
	return &repo
}

func (r *genreRepository) FindByNames(ctx context.Context, names []string) ([]*model.Genre, error) {
	g := query.Use(r.ReaderDB).Genre

	genres, err := g.WithContext(ctx).
		Where(g.Name.In(names...)).
		Find()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return make([]*model.Genre, 0), nil
	}

	return genres, err
}
