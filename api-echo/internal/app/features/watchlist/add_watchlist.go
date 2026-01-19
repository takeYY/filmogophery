package watchlist

import (
	"context"
	"fmt"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	AddWatchlistUseCase interface {
		Run(ctx context.Context, operator *model.Users, movieID int32, priority int32) error
	}

	addWatchlistInteractor struct {
		movieRepo     repositories.IMovieRepository
		watchlistRepo repositories.IWatchlistRepository
	}
)

func NewAddWatchlistInteractor(
	movieRepo repositories.IMovieRepository,
	watchlistRepo repositories.IWatchlistRepository,
) AddWatchlistUseCase {
	return &addWatchlistInteractor{
		movieRepo,
		watchlistRepo,
	}
}

func (i *addWatchlistInteractor) Run(ctx context.Context, operator *model.Users, movieID int32, priority int32) error {
	logger := logger.GetLogger()

	// 映画の存在確認
	movie, err := i.movieRepo.FindByID(ctx, movieID)
	if err != nil {
		logger.Error().Msgf("failed to fetch movie: %s", err.Error())
		return responses.InternalServerError()
	}
	if movie == nil {
		return responses.NotFoundError("movie", map[string][]string{"id": {fmt.Sprintf("%d", movieID)}})
	}

	// ウォッチリストに登録
	err = i.watchlistRepo.Create(ctx, nil, &model.Watchlist{
		UserID:   operator.ID,
		MovieID:  movie.ID,
		Priority: &priority,
	})
	if err != nil {
		logger.Error().Msgf("failed to create watchlist: %s", err.Error())
		return responses.InternalServerError()
	}

	return nil
}
