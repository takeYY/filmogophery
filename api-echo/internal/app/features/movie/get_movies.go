package movie

import (
	"context"

	"github.com/rs/zerolog"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetMoviesUseCase interface {
		Run(ctx context.Context, operator *model.Users, genre string, limit int32, offset int32) ([]types.Movie, error)
	}

	getMoviesInteractor struct {
		movieService services.IMovieService
	}
)

func NewGetMoviesInteractor(movieService services.IMovieService) GetMoviesUseCase {
	return &getMoviesInteractor{
		movieService,
	}
}

func (i *getMoviesInteractor) Run(ctx context.Context, operator *model.Users, genre string, limit int32, offset int32) ([]types.Movie, error) {
	log := zerolog.Ctx(ctx)

	// ユーザーがレビューした映画を取得（ジャンル絞り込み可）
	movies, err := i.movieService.GetReviewedMoviesByUser(ctx, operator.ID, genre, limit, offset)
	if err != nil {
		return nil, err
	}

	// レスポンスを作成
	response := make([]types.Movie, 0, len(movies))
	for _, m := range movies {
		response = append(response, types.Movie{
			ID:             m.ID,
			Title:          m.Title,
			Overview:       m.Overview,
			ReleaseDate:    types.ConvertTime2Date(m.ReleaseDate),
			RuntimeMinutes: m.RuntimeMinutes,
			PosterURL:      m.PosterURL,
			TmdbID:         m.TmdbID,
			Genres:         types.NewGenresByModel(m.Genres),
		})
	}
	log.Debug().Msg("successfully set movies response")

	return response, nil
}
