package search

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/internal/pkg/logger"
)

type (
	SearchMoviesUseCase interface {
		Run(ctx context.Context, title string, limit int32, offset int32) ([]types.Movie, error)
	}

	searchMoviesInteractor struct {
		db       *gorm.DB
		movieSvc services.IMovieService
		redisSvc services.IRedisService
		tmdbSvc  services.ITmdbService
	}
)

func NewSearchMoviesInteractor(
	db *gorm.DB,
	movieSvc services.IMovieService,
	redisSvc services.IRedisService,
	tmdbSvc services.ITmdbService,
) SearchMoviesUseCase {
	return &searchMoviesInteractor{
		db,
		movieSvc,
		redisSvc,
		tmdbSvc,
	}
}

func (i *searchMoviesInteractor) Run(ctx context.Context, title string, limit int32, offset int32) ([]types.Movie, error) {
	logger := logger.GetLogger()

	// redis 格納用のキャッシュキーを生成
	cacheKey := i.newCacheKey(title, limit, offset)

	// Redis から情報を取得（あれば）
	movies := i.getMoviesFromRedis(ctx, cacheKey)
	if movies != nil {
		return movies, nil
	}

	// Redis になければ TMDb APIから映画情報を取得
	tmdbMovies, err := i.tmdbSvc.GetMoviesByTitle(title, limit, offset)
	if err != nil {
		return nil, err
	}
	logger.Debug().Msg("successfully search movies")

	// 取得した映画のtmdbIDリストで既存映画を取得
	mvs, err := i.getExistingMovies(ctx, tmdbMovies.Results)
	if err != nil {
		return nil, err
	}
	logger.Debug().Msg("successfully get movies by tmdbID")

	// moviesテーブルにない映画情報を作成
	newMovies := i.newMoviesForCreation(tmdbMovies.Results, mvs)

	// moviesテーブルにない映画を一括登録
	createdMovies, err := i.batchCreateMovies(ctx, newMovies)
	if err != nil {
		return nil, err
	}

	// 既存と新規作成した映画をマージ
	var allMovies []*model.Movies
	if createdMovies == nil {
		allMovies = mvs
	} else {
		allMovies = append(mvs, createdMovies...)
	}

	// TMDBの順序に合わせて返却用の映画情報に変換
	resultMovies := i.newMoviesForResponse(tmdbMovies.Results, allMovies)

	// Redisにキャッシュ（24時間）
	if err := i.redisSvc.Set(ctx, cacheKey, resultMovies, 24*time.Hour); err != nil {
		logger.Warn().Err(err).Msg("failed to cache movies in redis")
	}

	return resultMovies, nil
}

func (i *searchMoviesInteractor) newCacheKey(title string, limit int32, offset int32) string {
	// タイトルを正規化し、limit/offsetを含めたキーを生成
	return fmt.Sprintf("movies:search:%s:limit:%d:offset:%d",
		strings.ToLower(strings.TrimSpace(title)), limit, offset)
}

func (i *searchMoviesInteractor) getMoviesFromRedis(ctx context.Context, key string) []types.Movie {
	logger := logger.GetLogger()

	// Redisからタイトルに一致する情報を取得
	var movies []types.Movie
	err := i.redisSvc.Get(ctx, key, &movies)
	if err == nil {
		// キャッシュヒット（limit/offset込みでキャッシュされているのでそのまま返す）
		logger.Debug().Msg("cache hit from redis")
		return movies
	}
	// redis.Nil以外のエラーはログ出力のみ（キャッシュ障害でサービス停止させない）
	if err != redis.Nil {
		logger.Warn().Err(err).Msg("redis get error")
	}

	return nil
}

func (i *searchMoviesInteractor) getExistingMovies(ctx context.Context, tmdbMovies []*types.TmdbMovieResult) ([]*model.Movies, error) {
	// 取得した映画のtmdbIDリストで既存映画を取得
	tmdbIDs := make([]int32, 0, len(tmdbMovies))
	for _, tmdbResult := range tmdbMovies {
		tmdbIDs = append(tmdbIDs, int32(tmdbResult.ID))
	}

	movies, err := i.movieSvc.GetMoviesByTmdbIDs(ctx, tmdbIDs)
	if err != nil {
		return nil, err
	}

	return movies, nil
}

func (i *searchMoviesInteractor) newMoviesForCreation(tmdbMovies []*types.TmdbMovieResult, existingMovies []*model.Movies) []*model.Movies {
	logger := logger.GetLogger()

	existingTmdbIDs := make(map[int32]bool)
	for _, mv := range existingMovies {
		existingTmdbIDs[mv.TmdbID] = true
	}

	newMovies := make([]*model.Movies, 0, len(tmdbMovies)-len(existingMovies))
	for _, tmdbMovie := range tmdbMovies {
		if existingTmdbIDs[int32(tmdbMovie.ID)] {
			continue
		}

		genres := make([]*model.Genres, 0, len(tmdbMovie.GenreIds))
		for _, g := range tmdbMovie.GenreIds {
			genres = append(genres, &model.Genres{ID: int32(*g)})
		}

		releaseDate, err := constant.ToTime(tmdbMovie.ReleaseDate)
		if err != nil {
			logger.Error().Msgf("failed to convert release_date to time")
			releaseDate = constant.GetDefaultDate()
		}

		newMovies = append(newMovies, &model.Movies{
			TmdbID:      int32(tmdbMovie.ID),
			Title:       tmdbMovie.Title,
			Overview:    tmdbMovie.Overview,
			ReleaseDate: releaseDate,
			PosterURL:   tmdbMovie.PosterPath,
			Genres:      genres,
		})
	}

	return newMovies
}

func (i *searchMoviesInteractor) batchCreateMovies(ctx context.Context, newMovies []*model.Movies) ([]*model.Movies, error) {
	if len(newMovies) == 0 {
		return nil, nil
	}

	logger := logger.GetLogger()

	err := i.db.Transaction(func(tx *gorm.DB) error {
		return i.movieSvc.BatchCreate(ctx, tx, newMovies)
	})
	if err != nil {
		logger.Error().Msgf("failed to batch create movies: %s", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}

	// 新規登録した映画のIDを取得
	newTmdbIDs := make([]int32, 0, len(newMovies))
	for _, mv := range newMovies {
		newTmdbIDs = append(newTmdbIDs, mv.TmdbID)
	}

	// ジャンル情報を含めて再取得
	newMoviesWithGenres, err := i.movieSvc.GetMoviesByTmdbIDs(ctx, newTmdbIDs)
	if err != nil {
		logger.Error().Msgf("failed to get newly created movies: %s", err.Error())
		return nil, echo.NewHTTPError(http.StatusInternalServerError, "system error")
	}

	return newMoviesWithGenres, nil
}

func (i *searchMoviesInteractor) newMoviesForResponse(tmdbMovies []*types.TmdbMovieResult, allMovies []*model.Movies) []types.Movie {
	// tmdbIDをキーにしたマップを作成
	movieMap := make(map[int32]*model.Movies)
	for _, mv := range allMovies {
		movieMap[mv.TmdbID] = mv
	}

	// TMDBの順序に従って結果を構築
	resultMovies := make([]types.Movie, 0, len(tmdbMovies))
	for _, tmdbMovie := range tmdbMovies {
		mv, ok := movieMap[int32(tmdbMovie.ID)]
		if !ok {
			continue
		}
		resultMovies = append(resultMovies, types.Movie{
			ID:             mv.ID,
			TmdbID:         mv.TmdbID,
			Title:          mv.Title,
			Overview:       mv.Overview,
			ReleaseDate:    constant.ToDate(mv.ReleaseDate),
			RuntimeMinutes: mv.RuntimeMinutes,
			PosterURL:      mv.PosterURL,
			Genres:         types.NewGenresByModel(mv.Genres),
		})
	}

	return resultMovies
}
