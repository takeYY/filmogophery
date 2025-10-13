package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IMovieService interface {
		// --- Read --- //

		// 映画を全て取得
		GetMovies(ctx context.Context, genre string, limit int32) ([]*model.Movies, error)
		// IDに一致する映画を取得
		GetMovieByID(ctx context.Context, movieID int32) (*model.Movies, error)
	}

	movieService struct {
		movieRepo repositories.IMovieRepository
	}
)

func NewMovieService(
	movieRepo repositories.IMovieRepository,
) IMovieService {
	return &movieService{
		movieRepo,
	}
}

func (s *movieService) GetMovies(ctx context.Context, genre string, limit int32) ([]*model.Movies, error) {
	logger := logger.GetLogger()

	movies, err := s.movieRepo.FindByGenre(ctx, genre, limit)
	if err != nil {
		logger.Error().Msgf("failed to get movies: %s", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}
	logger.Debug().Msg("successfully fetched movies")

	return movies, err
}

// IDに一致する映画を取得
func (s *movieService) GetMovieByID(ctx context.Context, movieID int32) (*model.Movies, error) {
	logger := logger.GetLogger()

	movie, err := s.movieRepo.FindByID(ctx, movieID)
	if err != nil {
		logger.Error().Msgf("failed to get a movie(id=%d): %s", movieID, err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}
	if movie == nil {
		em := fmt.Sprintf("movie(id=%d) is not found", movieID)
		logger.Error().Msg(em)
		return nil, echo.NewHTTPError(http.StatusNotFound, em)
	}
	logger.Debug().Msg("successfully fetched a movie")

	return movie, err
}
