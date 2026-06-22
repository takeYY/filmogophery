package services

import (
	"context"

	"github.com/rs/zerolog"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
)

type (
	IWatchHistoryService interface {
		// --- Create --- //

		// --- Read --- //

		// 映画IDに一致する視聴履歴を取得する
		GetByMovieID(ctx context.Context, operator *model.Users, movie *model.Movies) ([]*model.WatchHistory, error)
		// ユーザーIDに一致する視聴履歴を取得する
		GetByUserID(ctx context.Context, operator *model.Users, limit, offset int32) ([]*model.WatchHistory, error)

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
	log := zerolog.Ctx(ctx)

	whs, err := s.watchHistRepo.FindByMovieID(ctx, operator, movie)
	if err != nil {
		log.Error().Msgf("failed to get a watch history(movieID=%d): %s", movie.ID, err.Error())
		return nil, responses.InternalServerError()
	}

	return whs, nil
}

// ユーザーIDに一致する視聴履歴を取得する
func (s *watchHistoryService) GetByUserID(
	ctx context.Context, operator *model.Users, limit, offset int32,
) ([]*model.WatchHistory, error) {
	log := zerolog.Ctx(ctx)

	whs, err := s.watchHistRepo.FindByUserID(ctx, operator, limit, offset)
	if err != nil {
		log.Error().Msgf("failed to get a watch history(userID=%d): %s", operator.ID, err.Error())
		return nil, responses.InternalServerError()
	}

	return whs, nil
}
