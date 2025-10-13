package services

import (
	"gorm.io/gorm"

	"filmogophery/internal/app/repositories"
	"filmogophery/internal/config"
)

type (
	IServiceContainer interface {
		GenreService() IGenreService
		MovieService() IMovieService
		PlatformService() IPlatformService
		ReviewService() IReviewService
		TmdbService() ITmdbService
	}

	serviceContainer struct {
		db   *gorm.DB
		conf *config.Config
	}
)

func NewServiceContainer(db *gorm.DB, conf *config.Config) IServiceContainer {
	return &serviceContainer{
		db, conf,
	}
}

func (c *serviceContainer) GenreService() IGenreService {
	return NewGenreService(
		repositories.NewGenreRepository(c.db),
	)
}

func (c *serviceContainer) MovieService() IMovieService {
	return NewMovieService(
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
	)
}

func (c *serviceContainer) TmdbService() ITmdbService {
	return NewTmdbService(
		repositories.NewTmdbRepository(&c.conf.Tmdb),
	)
}
