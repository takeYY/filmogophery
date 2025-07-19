package services

import (
	"filmogophery/internal/app/repositories"

	"gorm.io/gorm"
)

type (
	IServiceContainer interface {
		MovieService() IMovieService
	}

	serviceContainer struct{ db *gorm.DB }
)

func NewServiceContainer(db *gorm.DB) IServiceContainer {
	return &serviceContainer{db: db}
}

func (c *serviceContainer) MovieService() IMovieService {
	return NewMovieService(
		repositories.NewMovieRepository(c.db),
	)
}
