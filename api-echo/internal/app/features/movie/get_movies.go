package movie

import (
	"context"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/logger"
)

type (
	GetMoviesUseCase interface {
		Run(ctx context.Context, genre string, limit int32) ([]types.Movie, error)
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

func (i *getMoviesInteractor) Run(ctx context.Context, genre string, limit int32) ([]types.Movie, error) {
	logger := logger.GetLogger()

	// 全ての映画を取得
	movies, err := i.movieService.GetMovies(ctx, genre, limit)
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
	logger.Debug().Msg("successfully set movies response")

	return response, nil
}
