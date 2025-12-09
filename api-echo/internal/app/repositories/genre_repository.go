package repositories

import (
	"context"

	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/gen/query"
)

type (
	IGenreRepository interface {
		// --- Create --- //

		// 映画ジャンルを紐づける
		BatchCreate(ctx context.Context, tx *gorm.DB, movieGenre []*model.MovieGenres) error

		// --- Read --- //

		// 全てのジャンルを取得
		FindAll(ctx context.Context) ([]*model.Genres, error)
		// 名前と一致するジャンルを取得
		FindByNames(ctx context.Context, names []string) ([]*model.Genres, error)
	}

	genreRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewGenreRepository(db *gorm.DB) IGenreRepository {
	return &genreRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// 映画ジャンルを紐づける
func (r *genreRepository) BatchCreate(ctx context.Context, tx *gorm.DB, movieGenre []*model.MovieGenres) error {
	mg := query.Use(r.WriterDB).MovieGenres
	if tx != nil {
		mg = query.Use(tx).MovieGenres
	}

	return mg.WithContext(ctx).Omit(field.AssociationFields).CreateInBatches(movieGenre, BATCH_SIZE)
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
