package platform

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
)

type (
	GetPlatformsUseCase interface {
		Run(ctx context.Context) ([]types.Platform, error)
	}

	getPlatformsInteractor struct {
		platformService services.IPlatformService
	}
)

func NewGetPlatformsInteractor(platformService services.IPlatformService) GetPlatformsUseCase {
	return &getPlatformsInteractor{
		platformService,
	}
}

func (i *getPlatformsInteractor) Run(ctx context.Context) ([]types.Platform, error) {
	// 全てのプラットフォームを取得
	platforms, err := i.platformService.GetAllPlatforms(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewPlatformsByModel(platforms), nil
}
