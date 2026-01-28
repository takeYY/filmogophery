package user

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetWatchHistoryUseCase interface {
		Run(ctx context.Context, operator *model.Users, limit int32, offset int32) ([]types.WatchHistory, error)
	}
	getWatchHistoryInteractor struct {
		watchHistSvc services.IWatchHistoryService
	}
)

func NewGetWatchHistoryInteractor(
	watchHistSvc services.IWatchHistoryService,
) GetWatchHistoryUseCase {
	return &getWatchHistoryInteractor{
		watchHistSvc,
	}
}

func (i *getWatchHistoryInteractor) Run(
	ctx context.Context, operator *model.Users, limit int32, offset int32,
) ([]types.WatchHistory, error) {
	// 視聴履歴を取得
	watchHistories, err := i.watchHistSvc.GetByUserID(ctx, operator, limit, offset)
	if err != nil {
		return nil, err
	}

	response := make([]types.WatchHistory, 0, len(watchHistories))
	for _, wh := range watchHistories {
		response = append(response, types.WatchHistory{
			ID:        wh.ID,
			WatchedAt: constant.ToDate(*wh.WatchedDate),
			Platform:  types.NewPlatformByModel(wh.Platform),
			Movie: types.Movie{
				ID:             wh.Movie.ID,
				Title:          wh.Movie.Title,
				Overview:       wh.Movie.Overview,
				ReleaseDate:    types.ConvertTime2Date(wh.Movie.ReleaseDate),
				RuntimeMinutes: wh.Movie.RuntimeMinutes,
				PosterURL:      wh.Movie.PosterURL,
				TmdbID:         wh.Movie.TmdbID,
				Genres:         types.NewGenresByModel(wh.Movie.Genres),
			},
		})
	}
	return response, nil
}
