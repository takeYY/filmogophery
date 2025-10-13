package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IReviewService interface {
		// IDに一致するレビューを取得
		GetReviewByID(ctx context.Context, userID int32, id int32) (*model.Reviews, error)
		// 映画IDに一致するレビューを取得
		GetReviewByMovieID(ctx context.Context, userID int32, movie *model.Movies) (*model.Reviews, error)

		// レビューIDに一致する視聴履歴を取得
		GetWatchHistoryByReviewID(ctx context.Context, review *model.Reviews) ([]*model.WatchHistory, error)
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

// IDに一致するレビューを取得
func (s *reviewService) GetReviewByID(ctx context.Context, userID int32, id int32) (*model.Reviews, error) {
	logger := logger.GetLogger()

	review, err := s.reviewRepo.FindByID(ctx, userID, id)
	if err != nil {
		logger.Error().Msgf("failed to get a review(userID=%d, id=%d): %s", userID, id, err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}
	if review == nil {
		em := fmt.Sprintf("review(id=%d) is not found", id)
		logger.Info().Msg(em)
		return nil, echo.NewHTTPError(http.StatusNotFound, em)
	}
	logger.Debug().Msg("successfully fetched a review")

	return review, err
}

// 映画IDに一致するレビューを取得
func (s *reviewService) GetReviewByMovieID(ctx context.Context, userID int32, movie *model.Movies) (*model.Reviews, error) {
	logger := logger.GetLogger()

	review, err := s.reviewRepo.FindByMovieID(ctx, userID, movie.ID)
	if err != nil {
		logger.Error().Msgf("failed to get a review(userID=%d, movieID=%d): %s", userID, movie.ID, err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}
	if review == nil {
		logger.Info().Msg("review is not found")
	}
	logger.Debug().Msg("successfully fetched a review")

	return review, err
}

// レビューIDに一致する視聴履歴を取得
func (s *reviewService) GetWatchHistoryByReviewID(ctx context.Context, review *model.Reviews) ([]*model.WatchHistory, error) {
	logger := logger.GetLogger()

	watchHistories, err := s.watchHistoryRepo.FindByReviewID(ctx, review.ID)
	if err != nil {
		logger.Error().Msgf("failed to get a watch history(reviewID=%d): %s", review.ID, err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}
	logger.Debug().Msg("successfully fetched watch histories")

	return watchHistories, err
}
