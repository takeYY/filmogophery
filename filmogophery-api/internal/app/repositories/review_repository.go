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
