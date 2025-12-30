package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/config"
)

type (
	IServiceContainer interface {
		DB() *gorm.DB
		GenreService() IGenreService
		MovieService() IMovieService
		PlatformService() IPlatformService
		ReviewService() IReviewService
		WatchHistoryService() IWatchHistoryService
		TmdbService() ITmdbService
		RedisService() IRedisService
	}

	serviceContainer struct {
		db    *gorm.DB
		conf  *config.Config
		redis *redis.Client
	}
)

func NewServiceContainer(db *gorm.DB, conf *config.Config, redisClient *redis.Client) IServiceContainer {
	return &serviceContainer{
		db, conf, redisClient,
	}
}

func (c *serviceContainer) DB() *gorm.DB {
	return c.db
}

func (c *serviceContainer) GenreService() IGenreService {
	return NewGenreService(
		repositories.NewGenreRepository(c.db),
	)
}

func (c *serviceContainer) MovieService() IMovieService {
	return NewMovieService(
		repositories.NewGenreRepository(c.db),
		repositories.NewMovieRepository(c.db),
	)
}

func (c *serviceContainer) PlatformService() IPlatformService {
	return NewPlatformService(
		repositories.NewPlatformRepository(c.db),
	)
}

func (c *serviceContainer) ReviewService() IReviewService {
	return NewReviewService(
		repositories.NewReviewRepository(c.db),
		repositories.NewWatchHistoryRepository(c.db),
	)
}

func (c *serviceContainer) WatchHistoryService() IWatchHistoryService {
	return NewWatchHistoryService(
		repositories.NewWatchHistoryRepository(c.db),
	)
}

func (c *serviceContainer) TmdbService() ITmdbService {
	return NewTmdbService(
		repositories.NewTmdbRepository(&c.conf.Tmdb),
	)
}

func (c *serviceContainer) RedisService() IRedisService {
	return NewRedisService(c.redis)
}
