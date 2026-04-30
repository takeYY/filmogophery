package user

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	GetUserPointsUseCase interface {
		Run(ctx context.Context, operator *model.Users, limit, offset int32) (*types.UserPoints, error)
	}

	getUserPointsInteractor struct {
		pointService services.IPointService
		pointRepo    repositories.IPointRepository
	}
)

func NewGetUserPointsInteractor(
	pointService services.IPointService,
	pointRepo repositories.IPointRepository,
) GetUserPointsUseCase {
	return &getUserPointsInteractor{
		pointService,
		pointRepo,
	}
}

func (i *getUserPointsInteractor) Run(ctx context.Context, operator *model.Users, limit, offset int32) (*types.UserPoints, error) {
	logger := logger.GetLogger()

	up, err := i.pointService.GetUserPoints(ctx, operator.ID)
	if err != nil {
		return nil, err
	}

	histories, err := i.pointRepo.FindHistoryByUserID(ctx, operator.ID, limit, offset)
	if err != nil {
		logger.Error().Msgf("failed to get point history(userID=%d): %s", operator.ID, err.Error())
		return nil, responses.InternalServerError()
	}

	items := make([]types.PointHistoryItem, 0, len(histories))
	for _, h := range histories {
		items = append(items, types.PointHistoryItem{
			ID:          h.ID,
			Points:      h.Points,
			Action:      h.Action,
			ReferenceID: h.ReferenceID,
			CreatedAt:   h.CreatedAt,
		})
	}

	nextLevelPoints := services.CalcNextLevelPoints(up.TotalPoints, *up.Level)
	currentLevelWidth := services.CalcCurrentLevelWidth(*up.Level)

	return &types.UserPoints{
		TotalPoints:       up.TotalPoints,
		Level:             *up.Level,
		NextLevelPoints:   nextLevelPoints,
		CurrentLevelWidth: currentLevelWidth,
		PointHistory:      items,
	}, nil
}
