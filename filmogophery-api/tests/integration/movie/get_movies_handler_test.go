package movie_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/tests/mocks"
)

// 正常系のテスト
func TestMakeGetMoviesHandler(t *testing.T) {
	// --- Arrange --- //
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	poster01 := "/1234567890.jpg"
	series01 := int32(314)
	genre011 := "ジャンル01_1"
	genre012 := "ジャンル01_2"

	movieRepo := mocks.NewMockIMovieRepository(ctrl)
	movieRepo.EXPECT().
		FindAll(gomock.All()).
		Return(
			[]*model.Movie{
				{
					ID:          1,
					Title:       "テスト01",
					Overview:    "テスト概要01",
					ReleaseDate: time.Date(2024, 12, 25, 0, 0, 0, 0, time.Local),
					RunTime:     123,
					PosterURL:   &poster01,
					SeriesID:    &series01,
					TmdbID:      1592,
					Genres: []*model.Genre{
						{
							ID:   65,
							Code: "genre01",
							Name: &genre011,
						},
						{
							ID:   35,
							Code: "genre02",
							Name: &genre012,
						},
					},
					Series:          nil,
					MovieImpression: nil,
				},
				{
					ID:              2,
					Title:           "テスト02",
					Overview:        "",
					ReleaseDate:     time.Date(2024, 12, 31, 0, 0, 0, 0, time.Local),
					RunTime:         987,
					PosterURL:       nil,
					SeriesID:        nil,
					TmdbID:          8979,
					Series:          nil,
					MovieImpression: nil,
				},
			},
			nil,
		)
	movieService := services.NewMovieService(movieRepo)
	handler := movie.BuildGetMoviesHandler(movieService)

	// --- Act --- //
	var actual []types.Movie
	var expected []types.Movie

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/movies", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/v1/movies")

	err := handler(c)

	// --- Expected --- //
	expectedJSON := `[
		{
			"id": 1,
			"title": "テスト01",
			"overview": "テスト概要01",
			"releaseDate": "2024-12-25 00:00:00 +0900 JST",
			"runTime": 123,
			"posterURL": "/1234567890.jpg",
			"tmdbID": 1592,
			"genres": [
				{
					"code": "genre01",
					"name": "ジャンル01_1"
				},
				{
					"code": "genre02",
					"name": "ジャンル01_2"
				}
			]
		},
		{
			"id": 2,
			"title": "テスト02",
			"overview": "",
			"releaseDate": "2024-12-31 00:00:00 +0900 JST",
			"runTime": 987,
			"posterURL": null,
			"tmdbID": 8979,
			"genres": []
		}
	]`

	// --- Assert --- //
	if !assert.NoError(t, err) {
		t.Errorf("handler returned an error: %v", err)
	}

	// Status code
	assert.Equal(t, http.StatusOK, rec.Code)

	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("expected unmarshal is failed")
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
		t.Fatalf("actual unmarshal is failed")
	}

	// JSON 形式で検証
	assert.ElementsMatch(t, expected, actual)
}

// 正常系のテスト
// データがない場合
func TestMakeGetMoviesHandler_WithoutData(t *testing.T) {
	// --- Arrange --- //
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	movieRepo := mocks.NewMockIMovieRepository(ctrl)
	movieRepo.EXPECT().
		FindAll(gomock.All()).
		Return(make([]*model.Movie, 0), nil)
	movieService := services.NewMovieService(movieRepo)
	handler := movie.BuildGetMoviesHandler(movieService)

	// --- Act --- //
	var actual []types.Movie
	var expected []types.Movie

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/v1/movies", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/v1/movies")

	err := handler(c)

	// --- Expected --- //
	expectedJSON := `[]`

	// --- Assert --- //
	if !assert.NoError(t, err) {
		t.Errorf("handler returned an error: %v", err)
	}

	// Status code
	assert.Equal(t, http.StatusOK, rec.Code)

	if err := json.Unmarshal([]byte(expectedJSON), &expected); err != nil {
		t.Fatalf("expected unmarshal is failed")
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &actual); err != nil {
		t.Fatalf("actual unmarshal is failed")
	}

	// JSON 形式で検証
	assert.ElementsMatch(t, expected, actual)
}
