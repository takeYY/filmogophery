package services

import (
	"context"
	"net/http"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"

	"github.com/labstack/echo/v4"
)

type (
	IReviewService interface {
		// IDに一致するレビューを取得
		GetReviewByMovieID(ctx context.Context, userID int32, movie *model.Movies) (*model.Reviews, error)
	}
	reviewService struct {
		reviewRepo repositories.IReviewRepository
	}
)

func NewReviewService(
	reviewRepo repositories.IReviewRepository,
) IReviewService {
	return &reviewService{
		reviewRepo,
	}
}

// IDに一致するレビューを取得
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
