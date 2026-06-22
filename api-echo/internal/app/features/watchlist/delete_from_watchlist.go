package watchlist

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
)

type (
	DeleteFromWatchlistUseCase interface {
		Run(ctx context.Context, operator *model.Users, watchlistID int32) error
	}

	deleteFromWatchlistInteractor struct {
		watchlistRepo repositories.IWatchlistRepository
	}
)

func NewDeleteFromWatchlistInteractor(
	watchlistRepo repositories.IWatchlistRepository,
) DeleteFromWatchlistUseCase {
	return &deleteFromWatchlistInteractor{
		watchlistRepo,
	}
}

func (i *deleteFromWatchlistInteractor) Run(ctx context.Context, operator *model.Users, watchlistID int32) error {
	log := zerolog.Ctx(ctx)

	// ウォッチリストから削除
	affected, err := i.watchlistRepo.DeleteByID(ctx, nil, watchlistID)
	if err != nil {
		log.Error().Msgf("failed to delete from watchlist(id=%d): %s", watchlistID, err.Error())
		return responses.InternalServerError()
	}
	if affected == 0 {
		return responses.NotFoundError("watchlist", map[string][]string{"id": {fmt.Sprintf("%d", watchlistID)}})
	}

	return nil
}
