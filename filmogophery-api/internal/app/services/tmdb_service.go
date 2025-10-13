package services

import (
	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/logger"
	"net/http"

	"github.com/labstack/echo/v4"
)

type (
	ITmdbService interface {
		// IDに一致する映画詳細を取得
		GetMovieDetailByID(id int32) (*types.TmdbMovieDetail, error)
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
