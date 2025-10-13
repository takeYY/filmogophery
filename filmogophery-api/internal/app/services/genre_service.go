package services

import (
	"context"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IGenreService interface {
		// --- Read --- //

		// 全てのジャンルを取得
		GetAllGenres(ctx context.Context) ([]*model.Genres, error)
	}

	genreService struct {
		genreRepo repositories.IGenreRepository
	}
)

func NewGenreService(
	genreRepo repositories.IGenreRepository,
) IGenreService {
	return &genreService{
		genreRepo,
	}
}

// 全てのジャンルを取得
func (s *genreService) GetAllGenres(ctx context.Context) ([]*model.Genres, error) {
	logger := logger.GetLogger()

	genres, err := s.genreRepo.FindAll(ctx)
	if err != nil {
		logger.Error().Msgf("failed to fetch genres: %s", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}

	return genres, nil
}
