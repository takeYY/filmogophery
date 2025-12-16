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
	IMovieRepository interface {
		// --- Create --- //

		// 映画を一括作成
		BatchCreate(ctx context.Context, tx *gorm.DB, movies []*model.Movies) error

		// --- Read --- //

		// ID に一致する映画を取得
		FindByID(ctx context.Context, id int32) (*model.Movies, error)
		// ジャンルを指定して取得
		FindByGenre(ctx context.Context, genre string, limit int32, offset int32) ([]*model.Movies, error)
		// tmdbIDs に一致する映画を取得
		FindByTmdbIDs(ctx context.Context, tmdbIDs []int32) ([]*model.Movies, error)
	}
	SaveMovieInput struct {
		Target *model.Movies
		Tx     *gorm.DB
	}

	movieRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

const (
	BATCH_SIZE = 100
)

func NewMovieRepository(db *gorm.DB) IMovieRepository {
	return &movieRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// 映画を一括作成
func (r *movieRepository) BatchCreate(ctx context.Context, tx *gorm.DB, movies []*model.Movies) error {
	m := query.Use(r.WriterDB).Movies
	if tx != nil {
		m = query.Use(tx).Movies
	}

	return m.WithContext(ctx).Omit(field.AssociationFields).CreateInBatches(movies, BATCH_SIZE)
}

// ID に一致する映画を取得
func (r *movieRepository) FindByID(ctx context.Context, id int32) (*model.Movies, error) {
	m := query.Use(r.ReaderDB).Movies

	result, err := m.WithContext(ctx).
		Preload(m.Genres).
		Preload(m.Series).
		Where(m.ID.Eq(id)).
		Take()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}

// ジャンルを指定して取得
func (r *movieRepository) FindByGenre(
	ctx context.Context, genre string, limit int32, offset int32,
) ([]*model.Movies, error) {
	m := query.Use(r.ReaderDB).Movies
	g := query.Use(r.ReaderDB).Genres
	mg := query.Use(r.ReaderDB).MovieGenres

	q := m.WithContext(ctx).Preload(m.Genres)
	if genre != "" {
		q = q.LeftJoin(mg, mg.MovieID.EqCol(m.ID)).
			LeftJoin(g, g.ID.EqCol(mg.GenreID)).
			Where(g.Code.Eq(genre))
	}

	return q.
		Order(m.CreatedAt.Desc()).
		Limit(int(limit)).
		Offset(int(offset)).
		Find()
}

// tmdbIDs に一致する映画を取得
func (r *movieRepository) FindByTmdbIDs(ctx context.Context, tmdbIDs []int32) ([]*model.Movies, error) {
	m := query.Use(r.ReaderDB).Movies

	return m.WithContext(ctx).
		Preload(m.Genres).
		Preload(m.Series).
		Where(m.TmdbID.In(tmdbIDs...)).
		Find()
}
