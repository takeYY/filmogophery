package services

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IPlatformService interface {
		// --- Read --- //

		// 全てのプラットフォームを取得
		GetAllPlatforms(ctx context.Context) ([]*model.Platforms, error)
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
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}

	return platforms, nil
}
