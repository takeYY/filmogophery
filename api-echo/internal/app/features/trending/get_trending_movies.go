package trending

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
)

type (
	GetTrendingMoviesUseCase interface {
		Run(ctx context.Context, operator *model.Users) ([]types.TrendingMovie, error)
	}

	getTrendingMoviesInteractor struct {
		db        *gorm.DB
		movieSvc  services.IMovieService
		reviewSvc services.IReviewService
		redisSvc  services.IRedisService
		tmdbSvc   services.ITmdbService
	}
)

func NewGetTrendingMoviesInteractor(
	db *gorm.DB,
	movieSvc services.IMovieService,
	reviewSvc services.IReviewService,
	redisSvc services.IRedisService,
	tmdbSvc services.ITmdbService,
) GetTrendingMoviesUseCase {
	return &getTrendingMoviesInteractor{
		db, movieSvc, reviewSvc, redisSvc, tmdbSvc,
	}
}

func (i *getTrendingMoviesInteractor) Run(ctx context.Context, operator *model.Users) ([]types.TrendingMovie, error) {
	log := zerolog.Ctx(ctx)

	// redis 格納用のキャッシュキーを生成
	cacheKey := i.newCacheKey()

	// Redis から情報を取得（あれば）
	movies := i.getMoviesFromRedis(ctx, cacheKey)
	if movies == nil {
		// Redis になければ TMDb APIから人気映画を取得
		tmdbMovies, err := i.tmdbSvc.GetTrendingMovies(ctx)
		if err != nil {
			return nil, err
		}
		log.Debug().Msg("successfully got Trending movies")

		// 取得した映画のtmdbIDリストで既存映画を取得
		mvs, err := i.getExistingMovies(ctx, tmdbMovies.Results)
		if err != nil {
			return nil, err
		}
		log.Debug().Msg("successfully got movies by tmdbID")

		// moviesテーブルにない映画情報を作成
		newMovies := i.newMoviesForCreation(ctx, tmdbMovies.Results, mvs)

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
		movies = i.newMoviesForResponse(tmdbMovies.Results, allMovies)

		// Redisにキャッシュ（24時間）
		if err := i.redisSvc.Set(ctx, cacheKey, movies, 24*time.Hour); err != nil {
			log.Warn().Err(err).Msg("failed to cache movies in redis")
		}
	}

	// レビュー済みフラグを付与
	movieIDs := make([]int32, 0, len(movies))
	for _, m := range movies {
		movieIDs = append(movieIDs, m.ID)
	}
	reviewedIDs, err := i.reviewSvc.GetReviewedMovieIDs(ctx, operator.ID, movieIDs)
	if err != nil {
		return nil, err
	}
	for idx := range movies {
		movies[idx].HasReview = reviewedIDs[movies[idx].ID]
	}

	return movies, nil
}

func (i *getTrendingMoviesInteractor) newCacheKey() string {
	return fmt.Sprintf("movies:trending:%s", time.Now().Format(constant.DateFormat))
}

func (i *getTrendingMoviesInteractor) getMoviesFromRedis(ctx context.Context, key string) []types.TrendingMovie {
	log := zerolog.Ctx(ctx)

	// Redisからタイトルに一致する情報を取得
	var movies []types.TrendingMovie
	err := i.redisSvc.Get(ctx, key, &movies)
	if err == nil {
		// キャッシュヒット
		log.Debug().Msg("cache hit from redis")
		return movies
	}
	// redis.Nil以外のエラーはログ出力のみ（キャッシュ障害でサービス停止させない）
	if err != redis.Nil {
		log.Warn().Err(err).Msg("redis get error")
	}

	return nil
}

func (i *getTrendingMoviesInteractor) getExistingMovies(ctx context.Context, tmdbMovies []types.TmdbTrendingMovie) ([]*model.Movies, error) {
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

func (i *getTrendingMoviesInteractor) newMoviesForCreation(ctx context.Context, tmdbMovies []types.TmdbTrendingMovie, existingMovies []*model.Movies) []*model.Movies {
	log := zerolog.Ctx(ctx)

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
			genres = append(genres, &model.Genres{ID: int32(g)})
		}

		releaseDate, err := constant.ToTime(tmdbMovie.ReleaseDate)
		if err != nil {
			log.Error().Msgf("failed to convert release_date to time")
			releaseDate = constant.GetDefaultDate()
		}

		newMovies = append(newMovies, &model.Movies{
			ID:             int32(tmdbMovie.ID),
			TmdbID:         int32(tmdbMovie.ID),
			Title:          tmdbMovie.Title,
			Overview:       tmdbMovie.Overview,
			ReleaseDate:    releaseDate,
			RuntimeMinutes: constant.DEFAULT_RUNTIME_MINUTES,
			PosterURL:      &tmdbMovie.PosterPath,
			Genres:         genres,
		})
	}

	return newMovies
}

func (i *getTrendingMoviesInteractor) batchCreateMovies(ctx context.Context, newMovies []*model.Movies) ([]*model.Movies, error) {
	if len(newMovies) == 0 {
		return nil, nil
	}

	log := zerolog.Ctx(ctx)

	err := i.db.Transaction(func(tx *gorm.DB) error {
		return i.movieSvc.BatchCreate(ctx, tx, newMovies)
	})
	if err != nil {
		log.Error().Msgf("failed to batch create movies: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	// 新規登録した映画のIDを取得
	newTmdbIDs := make([]int32, 0, len(newMovies))
	for _, mv := range newMovies {
		newTmdbIDs = append(newTmdbIDs, mv.TmdbID)
	}

	// ジャンル情報を含めて再取得
	newMoviesWithGenres, err := i.movieSvc.GetMoviesByTmdbIDs(ctx, newTmdbIDs)
	if err != nil {
		log.Error().Msgf("failed to get newly created movies: %s", err.Error())
		return nil, responses.InternalServerError()
	}

	return newMoviesWithGenres, nil
}

func (i *getTrendingMoviesInteractor) newMoviesForResponse(tmdbMovies []types.TmdbTrendingMovie, allMovies []*model.Movies) []types.TrendingMovie {
	// tmdbIDをキーにしたマップを作成
	movieMap := make(map[int32]*model.Movies)
	for _, mv := range allMovies {
		movieMap[mv.TmdbID] = mv
	}

	// TMDBの順序に従って結果を構築
	resultMovies := make([]types.TrendingMovie, 0, len(tmdbMovies))
	for _, tmdbMovie := range tmdbMovies {
		mv, ok := movieMap[int32(tmdbMovie.ID)]
		if !ok {
			continue
		}
		resultMovies = append(resultMovies, types.TrendingMovie{
			ID:        mv.ID,
			TmdbID:    mv.TmdbID,
			Title:     mv.Title,
			PosterURL: mv.PosterURL,
		})
	}

	return resultMovies
}
