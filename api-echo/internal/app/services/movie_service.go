package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/pkg/gen/model"
)

type (
	IMovieService interface {
		// --- Create --- //

		// 映画を一括作成
		BatchCreate(ctx context.Context, tx *gorm.DB, movies []*model.Movies) error

		// --- Read --- //

		// 映画一覧を取得
		GetMovies(ctx context.Context, genre string, limit int32, offset int32) ([]*model.Movies, error)
		// ユーザーがレビューした映画一覧を取得（ジャンル絞り込み可）
		GetReviewedMoviesByUser(ctx context.Context, userID int32, genre string, limit int32, offset int32) ([]*model.Movies, error)
		// IDに一致する映画を取得
		GetMovieByID(ctx context.Context, movieID int32) (*model.Movies, error)
		// tmdbIDsに一致する映画を取得
		GetMoviesByTmdbIDs(ctx context.Context, tmdbIDs []int32) ([]*model.Movies, error)

		// --- Update --- //

		// 上映時間を更新
		UpdateRuntimeMinutes(ctx context.Context, tx *gorm.DB, movie *model.Movies) error
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
	log := zerolog.Ctx(ctx)

	// 映画を一括作成
	err := s.movieRepo.BatchCreate(ctx, tx, movies)
	if err != nil {
		log.Error().Msgf("failed to batch create movies: %s", err.Error())
		return responses.InternalServerError()
	}

	// 映画に紐付けるジャンルIDを収集
	genreIDSet := make(map[int32]bool)
	for _, m := range movies {
		for _, g := range m.Genres {
			genreIDSet[g.ID] = true
		}
	}

	// DBに存在するジャンルIDのみに絞り込む（外部キー制約違反を防ぐ）
	if len(genreIDSet) > 0 {
		allGenreIDs := make([]int32, 0, len(genreIDSet))
		for id := range genreIDSet {
			allGenreIDs = append(allGenreIDs, id)
		}
		existingGenres, err := s.genreRepo.FindByIDs(ctx, allGenreIDs)
		if err != nil {
			log.Error().Msgf("failed to fetch genres: %s", err.Error())
			return responses.InternalServerError()
		}
		existingGenreIDs := make(map[int32]bool, len(existingGenres))
		for _, g := range existingGenres {
			existingGenreIDs[g.ID] = true
		}

		// 存在するジャンルIDのみで movie_genres を構築
		movieGenres := make([]*model.MovieGenres, 0)
		for _, m := range movies {
			for _, g := range m.Genres {
				if existingGenreIDs[g.ID] {
					movieGenres = append(movieGenres, &model.MovieGenres{
						MovieID: m.ID,
						GenreID: g.ID,
					})
				}
			}
		}

		if len(movieGenres) > 0 {
			err = s.genreRepo.BatchCreate(ctx, tx, movieGenres)
			if err != nil {
				log.Error().Msgf("failed to batch create movie_genres: %s", err.Error())
				return responses.InternalServerError()
			}
		}
	}

	return nil
}

// 映画一覧を取得
func (s *movieService) GetMovies(ctx context.Context, genre string, limit int32, offset int32) ([]*model.Movies, error) {
	log := zerolog.Ctx(ctx)

	movies, err := s.movieRepo.FindByGenre(ctx, genre, limit, offset)
	if err != nil {
		log.Error().Msgf("failed to get movies: %s", err.Error())
		return nil, responses.InternalServerError()
	}
	log.Debug().Msg("successfully fetched movies")

	return movies, err
}

// ユーザーがレビューした映画一覧を取得（ジャンル絞り込み可）
func (s *movieService) GetReviewedMoviesByUser(ctx context.Context, userID int32, genre string, limit int32, offset int32) ([]*model.Movies, error) {
	log := zerolog.Ctx(ctx)

	movies, err := s.movieRepo.FindReviewedByUser(ctx, userID, genre, limit, offset)
	if err != nil {
		log.Error().Msgf("failed to get reviewed movies for user(id=%d): %s", userID, err.Error())
		return nil, responses.InternalServerError()
	}
	log.Debug().Msg("successfully fetched reviewed movies")

	return movies, err
}

// IDに一致する映画を取得
func (s *movieService) GetMovieByID(ctx context.Context, movieID int32) (*model.Movies, error) {
	log := zerolog.Ctx(ctx)

	movie, err := s.movieRepo.FindByID(ctx, movieID)
	if err != nil {
		log.Error().Msgf("failed to get a movie(id=%d): %s", movieID, err.Error())
		return nil, responses.InternalServerError()
	}
	if movie == nil {
		em := fmt.Sprintf("movie(id=%d) is not found", movieID)
		log.Error().Msg(em)
		return nil, responses.NotFoundError("movie", map[string][]string{"id": {fmt.Sprintf("%d", movieID)}})
	}
	log.Debug().Msg("successfully fetched a movie")

	return movie, err
}

// tmdbIDsに一致する映画を取得
func (s *movieService) GetMoviesByTmdbIDs(ctx context.Context, tmdbIDs []int32) ([]*model.Movies, error) {
	log := zerolog.Ctx(ctx)

	movies, err := s.movieRepo.FindByTmdbIDs(ctx, tmdbIDs)
	if err != nil {
		log.Error().Msgf("failed to get movies: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	return movies, nil
}

// 上映時間を更新
func (s *movieService) UpdateRuntimeMinutes(ctx context.Context, tx *gorm.DB, movie *model.Movies) error {
	log := zerolog.Ctx(ctx)

	err := s.movieRepo.UpdateRuntimeMinutes(ctx, tx, movie)
	if err != nil {
		log.Error().Msgf("failed to update movies: %s", err.Error())
		return responses.InternalServerError()
	}

	return nil
}
