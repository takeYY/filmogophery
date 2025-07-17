package services

import (
	"context"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IMovieService interface {
		// --- Read --- //

		// 映画を全て取得
		GetMovies(ctx context.Context) ([]*model.Movie, error)
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

func (s *movieService) GetMovies(ctx context.Context) ([]*model.Movie, error) {
	logger := logger.GetLogger()

	movies, err := s.movieRepo.FindAll(ctx)
	if err != nil {
		logger.Error().Msgf("failed to get movies: %s", err.Error())
		return nil, err
	}
	logger.Debug().Msg("successfully fetch movies")

	return movies, err
}
