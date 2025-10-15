package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/tests"
)

// 正常系
// 映画データなし
func TestGetMoviesHandler_handle_WithoutDate(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	// サービス作成（トランザクション使用）
	svc := services.NewServiceContainer(tx, conf)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &getMoviesHandler{
		interactor: movie.NewGetMoviesInteractor(svc.MovieService()),
	}

	// Act
	err := handler.handle(c)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, `[]`, rec.Body.String())
}

// 正常系
// 映画データあり
func TestGetMoviesHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	// テストデータ作成
	tests.CreateMovies(t, tx, &model.Movies{ // ポスターとジャンルなし
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 11, 0, 0, 10, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})
	m2 := tests.CreateMovies(t, tx, &model.Movies{ // ポスターもジャンルもあり
		ID:             2,
		Title:          "テスト映画タイトル2",
		Overview:       "テスト概要2",
		ReleaseDate:    time.Date(2025, 10, 12, 0, 0, 20, 456789, time.Local),
		PosterURL:      &[]string{"/example.jpg"}[0],
		RuntimeMinutes: 6535,
		TmdbID:         8979,
	})
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 1})  // アクション
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 15}) // SF

	// サービス作成（トランザクション使用）
	svc := services.NewServiceContainer(tx, conf)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &getMoviesHandler{
		interactor: movie.NewGetMoviesInteractor(svc.MovieService()),
	}

	// Act
	err := handler.handle(c)

	// Expected
	expected := `[
		{
			"id": 1,
			"title": "テスト映画タイトル1",
			"overview": "テスト概要1",
			"releaseDate": "2025-10-11",
			"runtimeMinutes": 314,
			"posterURL": null,
			"tmdbID": 1592,
			"genres": []
		},
		{
			"id": 2,
			"title": "テスト映画タイトル2",
			"overview": "テスト概要2",
			"releaseDate": "2025-10-12",
			"runtimeMinutes": 6535,
			"posterURL": "/example.jpg",
			"tmdbID": 8979,
			"genres": [
				{"code": "action", "name": "アクション"},
				{"code": "sf", "name": "SF"}
			]
		}
	]`

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
