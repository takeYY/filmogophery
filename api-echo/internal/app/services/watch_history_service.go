package services

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IWatchHistoryService interface {
		// --- Create --- //

		// --- Read --- //

		// 映画IDに一致する視聴履歴を取得する
		GetByMovieID(ctx context.Context, operator *model.Users, movie *model.Movies) ([]*model.WatchHistory, error)

		// --- Update --- //

		// --- Delete --- //
	}
	watchHistoryService struct {
		watchHistRepo repositories.IWatchHistoryRepository
	}
)

func NewWatchHistoryService(
	watchHistRepo repositories.IWatchHistoryRepository,
) IWatchHistoryService {
	return &watchHistoryService{
		watchHistRepo,
	}
}

// 映画IDに一致する視聴履歴を取得する
func (s *watchHistoryService) GetByMovieID(
	ctx context.Context, operator *model.Users, movie *model.Movies,
) ([]*model.WatchHistory, error) {
	logger := logger.GetLogger()

	whs, err := s.watchHistRepo.FindByMovieID(ctx, operator, movie)
	if err != nil {
		logger.Error().Msgf("failed to get a watch history(movieID=%d): %s", movie.ID, err.Error())
		return nil, responses.InternalServerError()
	}

	return whs, nil
}
