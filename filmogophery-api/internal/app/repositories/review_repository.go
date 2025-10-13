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
	IReviewRepository interface {
		// --- Create --- //

		// レビューを作成
		Save(ctx context.Context, tx *gorm.DB, review *model.Reviews) error

		// IDに一致するレビューを取得
		FindByID(ctx context.Context, userID int32, id int32) (*model.Reviews, error)
		// 映画IDに一致するレビューを取得
		FindByMovieID(ctx context.Context, userID int32, movieID int32) (*model.Reviews, error)
	}
	reviewRepository struct {
		ReaderDB *gorm.DB
		WriterDB *gorm.DB
	}
)

func NewReviewRepository(db *gorm.DB) IReviewRepository {
	return &reviewRepository{
		ReaderDB: db.Clauses(dbresolver.Read),
		WriterDB: db.Clauses(dbresolver.Write),
	}
}

// レビューを作成
func (r *reviewRepository) Save(ctx context.Context, tx *gorm.DB, review *model.Reviews) error {
	rv := query.Use(r.WriterDB).Reviews
	if tx != nil {
		rv = query.Use(tx).Reviews
	}

	return rv.WithContext(ctx).Create(review)
}

// IDに一致するレビューを取得
func (r *reviewRepository) FindByID(ctx context.Context, userID int32, id int32) (*model.Reviews, error) {
	rv := query.Use(r.ReaderDB).Reviews

	result, err := rv.WithContext(ctx).
		Where(
			rv.ID.Eq(id),
			rv.UserID.Eq(userID),
		).
		Take()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}

// 映画IDに一致するレビューを取得
func (r *reviewRepository) FindByMovieID(ctx context.Context, userID int32, movieID int32) (*model.Reviews, error) {
	rv := query.Use(r.ReaderDB).Reviews

	result, err := rv.WithContext(ctx).
		Where(
			rv.UserID.Eq(userID),
			rv.MovieID.Eq(movieID),
		).
		Take()
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return result, err
}
