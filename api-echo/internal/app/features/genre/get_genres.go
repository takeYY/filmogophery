package genre

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
)

type (
	GetGenresUseCase interface {
		Run(ctx context.Context) ([]types.Genre, error)
	}

	getGenresInteractor struct {
		genreService services.IGenreService
	}
)

func NewGetGenresInteractor(genreService services.IGenreService) GetGenresUseCase {
	return &getGenresInteractor{
		genreService,
	}
}

func (i *getGenresInteractor) Run(ctx context.Context) ([]types.Genre, error) {
	// 全てのジャンルを取得
	genres, err := i.genreService.GetAllGenres(ctx)
	if err != nil {
		return nil, err
	}

	return types.NewGenresByModel(genres), nil
}
