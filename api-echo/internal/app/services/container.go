package services

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/hasher"
	"filmogophery/internal/pkg/jwt"
)

type (
	IServiceContainer interface {
		DB() *gorm.DB
		GenreService() IGenreService
		MovieService() IMovieService
		PlatformService() IPlatformService
		ReviewService() IReviewService
		UserService() IUserService
		WatchHistoryService() IWatchHistoryService
		TmdbService() ITmdbService
		RedisService() IRedisService
	}

	serviceContainer struct {
		db     *gorm.DB
		conf   *config.Config
		hasher *hasher.IPasswordHasher
		redis  *redis.Client
		token  *jwt.ITokenGenerator
	}
)

func NewServiceContainer(
	db *gorm.DB,
	conf *config.Config,
	hasher *hasher.IPasswordHasher,
	redisClient *redis.Client,
	token *jwt.ITokenGenerator,
) IServiceContainer {
	return &serviceContainer{
		db,
		conf,
		hasher,
		redisClient,
		token,
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

func (c *serviceContainer) UserService() IUserService {
	return NewUserService(
		c.db,
		NewAuthService(
			c.token,
			repositories.NewTokenRepository(c.db),
		),
		c.hasher,
		repositories.NewUserRepository(c.db),
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
