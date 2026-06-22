package services

import (
	"context"

	"github.com/rs/zerolog"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/types"
)

type (
	ITmdbService interface {
		// IDに一致する映画詳細を取得
		GetMovieDetailByID(ctx context.Context, id int32) (*types.TmdbMovieDetail, error)
		// タイトルに一致する映画を取得
		GetMoviesByTitle(ctx context.Context, title string, limit int32, offset int32) (*types.TmdbSearchMovieResult, error)
		// トレンド映画を取得
		GetTrendingMovies(ctx context.Context) (*types.TmdbTrendingMovieResult, error)
	}

	tmdbService struct {
		tmdbRepo repositories.ITmdbRepository
	}
)

func NewTmdbService(tmdbRepo repositories.ITmdbRepository) ITmdbService {
	return &tmdbService{
		tmdbRepo,
	}
}

// IDに一致する映画詳細を取得
func (s *tmdbService) GetMovieDetailByID(ctx context.Context, id int32) (*types.TmdbMovieDetail, error) {
	log := zerolog.Ctx(ctx)

	movieDetail, err := s.tmdbRepo.GetMovieDetail(id)
	if err != nil {
		log.Error().Msgf("failed to get a movie(id=%d) detail: %s", id, err.Error())
		return nil, responses.InternalServerError()
	}
	log.Debug().Msg("successfully fetch tmdb movie detail")

	return movieDetail, nil
}

// タイトルに一致する映画を取得
func (s *tmdbService) GetMoviesByTitle(ctx context.Context, title string, limit int32, offset int32) (*types.TmdbSearchMovieResult, error) {
	log := zerolog.Ctx(ctx)

	// pageを計算
	page := (offset / 20) + 1

	movies, err := s.tmdbRepo.GetMoviesByTitle(title, page)
	if err != nil {
		log.Error().Msgf("failed to fetch movies from tmdb: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	// ページ内のオフセット位置を計算
	pageOffset := offset % 20

	// 結果をlimitで切り出し
	if int32(len(movies.Results)) > pageOffset {
		end := pageOffset + limit
		if end > int32(len(movies.Results)) {
			end = int32(len(movies.Results))
		}
		movies.Results = movies.Results[pageOffset:end]
	} else {
		movies.Results = []*types.TmdbMovieResult{}
	}

	return movies, nil
}

// トレンド映画を取得
func (s *tmdbService) GetTrendingMovies(ctx context.Context) (*types.TmdbTrendingMovieResult, error) {
	log := zerolog.Ctx(ctx)

	movies, err := s.tmdbRepo.GetTrendingMovies()
	if err != nil {
		log.Error().Msgf("failed to fetch movies from tmdb: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	return movies, nil
}
