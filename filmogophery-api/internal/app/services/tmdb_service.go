package services

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/logger"
)

type (
	ITmdbService interface {
		// IDに一致する映画詳細を取得
		GetMovieDetailByID(id int32) (*types.TmdbMovieDetail, error)
		// タイトルに一致する映画を取得
		GetMoviesByTitle(title string, limit int32, offset int32) (*types.TmdbSearchMovieResult, error)
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
func (s *tmdbService) GetMovieDetailByID(id int32) (*types.TmdbMovieDetail, error) {
	logger := logger.GetLogger()

	movieDetail, err := s.tmdbRepo.GetMovieDetail(id)
	if err != nil {
		logger.Error().Msgf("failed to get a movie(id=%d) detail: %s", id, err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}
	logger.Debug().Msg("successfully fetch tmdb movie detail")

	return movieDetail, nil
}

// タイトルに一致する映画を取得
func (s *tmdbService) GetMoviesByTitle(title string, limit int32, offset int32) (*types.TmdbSearchMovieResult, error) {
	logger := logger.GetLogger()

	// pageを計算
	page := (offset / 20) + 1

	movies, err := s.tmdbRepo.GetMoviesByTitle(title, page)
	if err != nil {
		logger.Error().Msgf("failed to fetch movies from tmdb: %s", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
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
