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

		// 映画を作成
		Save(ctx context.Context, input SaveMovieInput) error

		// --- Read --- //

		// ID と一致する映画を取得
		FindByID(ctx context.Context, id *int32) (*model.Movie, error)
		// 全ての映画を取得
		FindAll(ctx context.Context) ([]*model.Movie, error)
	}
	SaveMovieInput struct {
		Target *model.Movie
		Tx     *gorm.DB
	}

	movieRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewMovieRepository(db *gorm.DB) *IMovieRepository {
	var repo IMovieRepository = &movieRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
	return &repo
}

func (r *movieRepository) FindByID(ctx context.Context, id *int32) (*model.Movie, error) {
	if id == nil {
		return nil, nil
	}

	m := query.Use(r.ReaderDB).Movie

	movie, err := m.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Where(m.ID.Eq(*id)).
		First()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return nil, nil
	}

	return movie, err
}

func (r *movieRepository) FindAll(ctx context.Context) ([]*model.Movie, error) {
	m := query.Use(r.ReaderDB).Movie

	movies, err := m.WithContext(ctx).
		Preload(field.Associations.Scopes(field.RelationFieldUnscoped)).
		Find()
	if errors.Is(err, gorm.ErrRecordNotFound) { // 0 件の場合
		return make([]*model.Movie, 0), nil
	}

	return movies, err
}

func (r *movieRepository) Save(ctx context.Context, input SaveMovieInput) error {
	m := query.Use(r.WriterDB).Movie
	if input.Tx != nil {
		m = query.Use(input.Tx).Movie
	}

	return m.WithContext(ctx).Create(input.Target)
}
