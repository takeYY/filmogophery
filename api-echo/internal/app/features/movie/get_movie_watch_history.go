package movie

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetMovieWatchHistoryUseCase interface {
		Run(ctx context.Context, operator *model.Users, movieID int32) ([]*types.MovieWatchHistory, error)
	}

	getMovieWatchHistoryInteractor struct {
		movieSvc     services.IMovieService
		watchHistSvc services.IWatchHistoryService
	}
)

func NewGetMovieWatchHistoryInteractor(
	movieSvc services.IMovieService,
	watchHistSvc services.IWatchHistoryService,
) GetMovieWatchHistoryUseCase {
	return &getMovieWatchHistoryInteractor{
		movieSvc,
		watchHistSvc,
	}
}

func (i *getMovieWatchHistoryInteractor) Run(
	ctx context.Context, operator *model.Users, movieID int32,
) ([]*types.MovieWatchHistory, error) {
	// 映画の存在確認
	movie, err := i.movieSvc.GetMovieByID(ctx, movieID)
	if err != nil {
		return nil, err
	}

	// 視聴履歴を取得
	watchHistories, err := i.watchHistSvc.GetByMovieID(ctx, operator, movie)
	if err != nil {
		return nil, err
	}

	response := make([]*types.MovieWatchHistory, 0, len(watchHistories))
	for _, wh := range watchHistories {
		response = append(response, &types.MovieWatchHistory{
			ID:        wh.ID,
			Platform:  types.NewPlatformByModel(wh.Platform),
			WatchedAt: *wh.WatchedDate,
		})
	}

	return response, nil
}
