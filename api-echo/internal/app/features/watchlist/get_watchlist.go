package watchlist

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	GetWatchlistUseCase interface {
		Run(ctx context.Context, operator *model.Users, limit int32, offset int32) ([]types.Watchlist, error)
	}

	getWatchlistInteractor struct {
		watchlistRepo repositories.IWatchlistRepository
	}
)

func NewGetWatchlistInteractor(
	watchlistRepo repositories.IWatchlistRepository,
) GetWatchlistUseCase {
	return &getWatchlistInteractor{
		watchlistRepo,
	}
}

func (i *getWatchlistInteractor) Run(
	ctx context.Context, operator *model.Users, limit int32, offset int32,
) ([]types.Watchlist, error) {
	logger := logger.GetLogger()

	// ウォッチリストを取得
	watchlist, err := i.watchlistRepo.FindByUserID(ctx, operator, limit, offset)
	if err != nil {
		logger.Error().Msg("failed to fetch watchlist")
		return nil, responses.InternalServerError()
	}

	// レスポンスを作成
	response := make([]types.Watchlist, 0, len(watchlist))
	for _, wl := range watchlist {
		m := wl.Movie
		response = append(response, types.Watchlist{
			ID:       wl.ID,
			AddedAt:  *wl.AddedAt,
			Priority: *wl.Priority,
			Movie: types.Movie{
				ID:             m.ID,
				Title:          m.Title,
				Overview:       m.Overview,
				ReleaseDate:    types.ConvertTime2Date(m.ReleaseDate),
				RuntimeMinutes: m.RuntimeMinutes,
				PosterURL:      m.PosterURL,
				TmdbID:         m.TmdbID,
				Genres:         types.NewGenresByModel(m.Genres),
			},
		})
	}

	return response, nil
}
