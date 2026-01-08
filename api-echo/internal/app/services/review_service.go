package services

import (
	"context"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IReviewService interface {
		// --- Create --- //

		// レビューを作成
		CreateReview(ctx context.Context, tx *gorm.DB, operator *model.Users, movie *model.Movies, rating *float64, comment *string) error
		// 視聴履歴を作成
		CreateWatchHistory(ctx context.Context, tx *gorm.DB, operator *model.Users, review *model.Reviews, platform *model.Platforms, watchedDate *constant.Date) error

		// --- Read --- //

		// IDに一致するレビューを取得
		GetReviewByID(ctx context.Context, operator *model.Users, id int32) (*model.Reviews, error)
		// 映画IDに一致するレビューを取得
		GetReviewByMovieID(ctx context.Context, operator *model.Users, movie *model.Movies) (*model.Reviews, error)

		// レビューIDに一致する視聴履歴を取得
		GetWatchHistoryByReviewID(ctx context.Context, operator *model.Users, review *model.Reviews) ([]*model.WatchHistory, error)

		// --- Update -- //

		// レビューを更新
		UpdateReview(ctx context.Context, tx *gorm.DB, review *model.Reviews, rating *float64, comment *string) error
	}
	reviewService struct {
		reviewRepo       repositories.IReviewRepository
		watchHistoryRepo repositories.IWatchHistoryRepository
	}
)

func NewReviewService(
	reviewRepo repositories.IReviewRepository,
	watchHistoryRepo repositories.IWatchHistoryRepository,
) IReviewService {
	return &reviewService{
		reviewRepo,
		watchHistoryRepo,
	}
}

// レビューを作成
func (s *reviewService) CreateReview(
	ctx context.Context, tx *gorm.DB, operator *model.Users, movie *model.Movies, rating *float64, comment *string,
) error {
	logger := logger.GetLogger()

	err := s.reviewRepo.Save(
		ctx,
		tx,
		&model.Reviews{
			UserID:  operator.ID,
			MovieID: movie.ID,
			Rating:  rating,
			Comment: comment,
		},
	)
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		em := fmt.Sprintf("review already exists for this movie(id=%d): %s", movie.ID, err.Error())
		logger.Error().Msg(em)
		return responses.BadRequestError(map[string][]string{"review": {"already exists"}})
	}
	if err != nil {
		logger.Error().Msgf("failed to create review: %s", err.Error())
		return responses.InternalServerError()
	}

	return nil
}

// 視聴履歴を作成
func (s *reviewService) CreateWatchHistory(
	ctx context.Context,
	tx *gorm.DB,
	operator *model.Users,
	review *model.Reviews,
	platform *model.Platforms,
	watchedDate *constant.Date,
) error {
	logger := logger.GetLogger()

	var parsedWatchedDate *time.Time
	if watchedDate != nil {
		parsedDate, err := time.ParseInLocation(constant.DateFormat, string(*watchedDate), time.Local)
		if err != nil {
			em := fmt.Sprintf("failed to parse watchedDate: %s", err.Error())
			logger.Error().Msg(em)
			return responses.BadRequestError(map[string][]string{"WatchedDate": {"failed to parse date"}})
		}
		parsedWatchedDate = &parsedDate
	}

	watchHistory := &model.WatchHistory{
		UserID:      operator.ID,
		MovieID:     review.MovieID,
		PlatformID:  platform.ID,
		WatchedDate: parsedWatchedDate,
	}
	err := s.watchHistoryRepo.Save(ctx, tx, watchHistory)
	if err != nil {
		logger.Error().Msgf("failed to create watch_history: %s", err.Error())
		return responses.InternalServerError()
	}
	logger.Debug().Msg("successfully created watch history")

	return nil
}

// IDに一致するレビューを取得
func (s *reviewService) GetReviewByID(ctx context.Context, operator *model.Users, id int32) (*model.Reviews, error) {
	logger := logger.GetLogger()

	review, err := s.reviewRepo.FindByID(ctx, operator.ID, id)
	if err != nil {
		logger.Error().Msgf("failed to get a review(userID=%d, id=%d): %s", operator.ID, id, err.Error())
		return nil, responses.InternalServerError()
	}
	if review == nil {
		em := fmt.Sprintf("review(id=%d) is not found", id)
		logger.Info().Msg(em)
		return nil, responses.NotFoundError("review", map[string][]string{"id": {fmt.Sprintf("%d", id)}})
	}
	logger.Debug().Msg("successfully fetched a review")

	return review, err
}

// 映画IDに一致するレビューを取得
func (s *reviewService) GetReviewByMovieID(ctx context.Context, operator *model.Users, movie *model.Movies) (*model.Reviews, error) {
	logger := logger.GetLogger()

	review, err := s.reviewRepo.FindByMovieID(ctx, operator.ID, movie.ID)
	if err != nil {
		logger.Error().Msgf("failed to get a review(userID=%d, movieID=%d): %s", operator.ID, movie.ID, err.Error())
		return nil, responses.InternalServerError()
	}
	if review == nil {
		logger.Info().Msg("review is not found")
	}
	logger.Debug().Msg("successfully fetched a review")

	return review, err
}

// レビューIDに一致する視聴履歴を取得
func (s *reviewService) GetWatchHistoryByReviewID(
	ctx context.Context, operator *model.Users, review *model.Reviews,
) ([]*model.WatchHistory, error) {
	logger := logger.GetLogger()

	watchHistories, err := s.watchHistoryRepo.FindByMovieID(ctx, operator, &review.Movie)
	if err != nil {
		logger.Error().Msgf("failed to get a watch history(reviewID=%d): %s", review.ID, err.Error())
		return nil, responses.InternalServerError()
	}
	logger.Debug().Msg("successfully fetched watch histories")

	return watchHistories, err
}

// レビューを更新
func (s *reviewService) UpdateReview(ctx context.Context, tx *gorm.DB, review *model.Reviews, rating *float64, comment *string) error {
	logger := logger.GetLogger()

	review.Rating = rating
	review.Comment = comment
	err := s.reviewRepo.Update(ctx, tx, review)
	if err != nil {
		logger.Error().Msgf("failed to update review(id=%d): %s", review.ID, err.Error())
		return responses.InternalServerError()
	}

	return nil
}
