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
// レビューデータなし（空配列を返す）
func TestGetMoviesHandler_handle_WithoutData(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	svc := services.NewServiceContainer(tx, conf, nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("operator", &model.Users{ID: 1})

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
// レビュー済み映画データあり（ジャンル絞り込みなし）
func TestGetMoviesHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	m1 := tests.CreateMovies(t, tx, &model.Movies{
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 11, 0, 0, 10, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})
	m2 := tests.CreateMovies(t, tx, &model.Movies{
		ID:             2,
		Title:          "テスト映画タイトル2",
		Overview:       "テスト概要2",
		ReleaseDate:    time.Date(2025, 10, 12, 0, 0, 20, 456789, time.Local),
		PosterURL:      &[]string{"/example.jpg"}[0],
		RuntimeMinutes: 6535,
		TmdbID:         8979,
	})
	tests.CreateMovies(t, tx, &model.Movies{ // レビューなし（一覧に出ないはず）
		ID:             3,
		Title:          "レビューなし映画",
		Overview:       "概要3",
		ReleaseDate:    time.Date(2025, 10, 13, 0, 0, 0, 0, time.Local),
		RuntimeMinutes: 100,
		TmdbID:         9999,
	})
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 28})  // アクション
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 878}) // SF

	// ユーザー1のレビューを作成（m1, m2 のみ）
	tests.CreateReviews(t, tx, &model.Reviews{UserID: 1, MovieID: m1.ID})
	tests.CreateReviews(t, tx, &model.Reviews{UserID: 1, MovieID: m2.ID})

	svc := services.NewServiceContainer(tx, conf, nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movies", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("operator", &model.Users{ID: 1})

	handler := &getMoviesHandler{
		interactor: movie.NewGetMoviesInteractor(svc.MovieService()),
	}

	// Act
	err := handler.handle(c)

	// Expected（レビューした2件のみ返る）
	expected := `[
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
		},
		{
			"id": 1,
			"title": "テスト映画タイトル1",
			"overview": "テスト概要1",
			"releaseDate": "2025-10-11",
			"runtimeMinutes": 314,
			"posterURL": null,
			"tmdbID": 1592,
			"genres": []
		}
	]`

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}

// 正常系
// レビュー済み映画データあり（ジャンル絞り込みあり）
func TestGetMoviesHandler_handle_WithGenre(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	m1 := tests.CreateMovies(t, tx, &model.Movies{
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 11, 0, 0, 10, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})
	m2 := tests.CreateMovies(t, tx, &model.Movies{
		ID:             2,
		Title:          "テスト映画タイトル2",
		Overview:       "テスト概要2",
		ReleaseDate:    time.Date(2025, 10, 12, 0, 0, 20, 456789, time.Local),
		PosterURL:      &[]string{"/example.jpg"}[0],
		RuntimeMinutes: 6535,
		TmdbID:         8979,
	})
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 28})  // アクション
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 878}) // SF

	// 両方レビュー済み
	tests.CreateReviews(t, tx, &model.Reviews{UserID: 1, MovieID: m1.ID})
	tests.CreateReviews(t, tx, &model.Reviews{UserID: 1, MovieID: m2.ID})

	svc := services.NewServiceContainer(tx, conf, nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/movies?genre=action", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("operator", &model.Users{ID: 1})

	handler := &getMoviesHandler{
		interactor: movie.NewGetMoviesInteractor(svc.MovieService()),
	}

	// Act
	err := handler.handle(c)

	// Expected（アクションジャンルのレビュー済み映画のみ返る）
	expected := `[
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
