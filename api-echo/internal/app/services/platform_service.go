package services

import (
	"context"
	"fmt"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IPlatformService interface {
		// --- Read --- //

		// 全てのプラットフォームを取得
		GetAllPlatforms(ctx context.Context) ([]*model.Platforms, error)
		// IDに一致するプラットフォームを取得
		GetByID(ctx context.Context, id int32) (*model.Platforms, error)
	}

	platformService struct {
		platformRepo repositories.IPlatformRepository
	}
)

func NewPlatformService(
	platformRepo repositories.IPlatformRepository,
) IPlatformService {
	return &platformService{
		platformRepo,
	}
}

// 全てのプラットフォームを取得
func (s *platformService) GetAllPlatforms(ctx context.Context) ([]*model.Platforms, error) {
	logger := logger.GetLogger()

	platforms, err := s.platformRepo.FindAll(ctx)
	if err != nil {
		logger.Error().Msgf("failed to fetch platforms: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	return platforms, nil
}

// IDに一致するプラットフォームを取得
func (s *platformService) GetByID(ctx context.Context, id int32) (*model.Platforms, error) {
	logger := logger.GetLogger()

	platform, err := s.platformRepo.FindByID(ctx, id)
	if err != nil {
		logger.Error().Msgf("failed to fetch platforms: %s", err.Error())
		return nil, responses.InternalServerError()
	}
	if platform == nil {
		em := fmt.Sprintf("platform(id=%d) is not found", id)
		logger.Info().Msg(em)
		return nil, responses.NotFoundError("platform", map[string][]string{"id": {fmt.Sprintf("%d", id)}})
	}
	logger.Debug().Msg("successfully fetched a platform")

	return platform, nil
}
