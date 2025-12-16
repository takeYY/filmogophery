package services

import (
	"context"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	IMovieService interface {
		// --- Create --- //

		// 映画を一括作成
		BatchCreate(ctx context.Context, tx *gorm.DB, movies []*model.Movies) error

		// --- Read --- //

		// 映画一覧を取得
		GetMovies(ctx context.Context, genre string, limit int32, offset int32) ([]*model.Movies, error)
		// IDに一致する映画を取得
		GetMovieByID(ctx context.Context, movieID int32) (*model.Movies, error)
		// tmdbIDsに一致する映画を取得
		GetMoviesByTmdbIDs(ctx context.Context, tmdbIDs []int32) ([]*model.Movies, error)
	}

	movieService struct {
		genreRepo repositories.IGenreRepository
		movieRepo repositories.IMovieRepository
	}
)

func NewMovieService(
	genreRepo repositories.IGenreRepository,
	movieRepo repositories.IMovieRepository,
) IMovieService {
	return &movieService{
		genreRepo,
		movieRepo,
	}
}

// 映画を一括作成
func (s *movieService) BatchCreate(ctx context.Context, tx *gorm.DB, movies []*model.Movies) error {
	logger := logger.GetLogger()

	// 映画を一括作成
	err := s.movieRepo.BatchCreate(ctx, tx, movies)
	if err != nil {
		logger.Error().Msgf("failed to batch create movies: %s", err.Error())
		return echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}

	// 映画ジャンルを紐づける
	movieGenres := make([]*model.MovieGenres, 0)
	for _, m := range movies {
		for _, g := range m.Genres {
			movieGenres = append(movieGenres, &model.MovieGenres{
				MovieID: m.ID,
				GenreID: g.ID,
			})
		}
	}

	if len(movieGenres) > 0 {
		err = s.genreRepo.BatchCreate(ctx, tx, movieGenres)
		if err != nil {
			logger.Error().Msgf("failed to batch create movie_genres: %s", err.Error())
			return echo.NewHTTPError(http.StatusInternalServerError, "system error")
		}
	}

	return nil
}

// 映画一覧を取得
func (s *movieService) GetMovies(ctx context.Context, genre string, limit int32, offset int32) ([]*model.Movies, error) {
	logger := logger.GetLogger()

	movies, err := s.movieRepo.FindByGenre(ctx, genre, limit, offset)
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

// tmdbIDsに一致する映画を取得
func (s *movieService) GetMoviesByTmdbIDs(ctx context.Context, tmdbIDs []int32) ([]*model.Movies, error) {
	logger := logger.GetLogger()

	movies, err := s.movieRepo.FindByTmdbIDs(ctx, tmdbIDs)
	if err != nil {
		logger.Error().Msgf("failed to get movies: %s", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}

	return movies, nil
}
