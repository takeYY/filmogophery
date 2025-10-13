package tests

import (
	"testing"

	"gorm.io/gorm"

	"filmogophery/internal/pkg/gen/model"
)

func CreateMovies(t *testing.T, tx *gorm.DB, fixture *model.Movies) *model.Movies {
	if result := tx.Create(fixture); result.Error != nil {
		t.Errorf("failed to create movies: %s", result.Error.Error())
	}

	return fixture
}

func CreateMovieGenres(t *testing.T, tx *gorm.DB, fixture *model.MovieGenres) *model.MovieGenres {
	if result := tx.Create(fixture); result.Error != nil {
		t.Errorf("failed to create movies: %s", result.Error.Error())
	}

	return fixture
}
