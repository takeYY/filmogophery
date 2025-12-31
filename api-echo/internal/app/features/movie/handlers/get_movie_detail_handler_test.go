package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	"filmogophery/internal/app/features/movie"
	"filmogophery/internal/app/services"
	"filmogophery/internal/app/types"
	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/tests"
	"filmogophery/tests/mocks"
)

func TestGetMovieDetailHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// サービス作成（トランザクション使用）
	now := time.Date(2025, 10, 13, 10, 20, 30, 456789, time.Local)
	svc := services.NewServiceContainer(tx, conf, nil, nil, nil)

	m1 := tests.CreateMovies(t, tx, &model.Movies{ // ポスターとジャンルなし
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 11, 0, 0, 10, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})
	s := tests.CreateSeries(t, tx, &model.Series{
		Name: "テストシリーズポスター", PosterURL: &[]string{"/example2.jpg"}[0],
	})
	m2 := tests.CreateMovies(t, tx, &model.Movies{ // ポスターもジャンルもあり
		ID:             2,
		Title:          "テスト映画タイトル2",
		Overview:       "テスト概要2",
		ReleaseDate:    time.Date(2025, 10, 12, 0, 0, 20, 456789, time.Local),
		PosterURL:      &[]string{"/example.jpg"}[0],
		SeriesID:       &s.ID,
		RuntimeMinutes: 6535,
		TmdbID:         8979,
	})
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 28})  // アクション
	tests.CreateMovieGenres(t, tx, &model.MovieGenres{MovieID: m2.ID, GenreID: 878}) // SF
	tests.CreateReviews(t, tx, &model.Reviews{
		ID: 1, UserID: 1, MovieID: m2.ID, Rating: &[]float64{3.23}[0], Comment: &[]string{"テストコメント"}[0],
		CreatedAt: &now, UpdatedAt: &now,
	})

	tmdbSvc := mocks.NewMockITmdbService(ctrl)
	tmdbSvc.EXPECT().GetMovieDetailByID(gomock.Eq(m1.TmdbID)).Return(&types.TmdbMovieDetail{
		TmdbMovieCommon: types.TmdbMovieCommon{VoteAverage: 6.282, VoteCount: 4},
	}, nil)
	tmdbSvc.EXPECT().GetMovieDetailByID(gomock.Eq(m2.TmdbID)).Return(&types.TmdbMovieDetail{
		TmdbMovieCommon: types.TmdbMovieCommon{VoteAverage: 3.184, VoteCount: 92},
	}, nil)

	for _, tt := range []struct {
		testCase string
		movieID  int32
		expected string
	}{
		{
			testCase: "基本映画",
			movieID:  m1.ID,
			expected: `{
				"id": 1,
				"title": "テスト映画タイトル1",
				"overview": "テスト概要1",
				"releaseDate": "2025-10-11",
				"runtimeMinutes": 314,
				"posterURL": null,
				"tmdbID": 1592,
				"genres": [],
				"series": null,
				"review": null,
				"voteAverage": 3.1,
				"voteCount": 4
			}`,
		},
		{
			testCase: "紐付け情報が全て含まれる映画",
			movieID:  m2.ID,
			expected: `{
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
				],
				"series": {
					"name": "テストシリーズポスター",
					"posterURL": "/example2.jpg"
				},
				"review": {
					"id": 1,
					"rating": 3.2,
					"comment": "テストコメント",
					"createdAt": "2025-10-13T10:20:30+09:00",
					"updatedAt": "2025-10-13T10:20:30+09:00"
				},
				"voteAverage": 1.6,
				"voteCount": 92
			}`,
		},
		// TODO: レビューの評価がない場合のテストケース追加
		// TODO: レビューのコメントがない場合のテストケース追加
	} {
		t.Run(tt.testCase, func(t *testing.T) {
			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/movies/%d", tt.movieID), nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(strconv.Itoa(int(tt.movieID)))

			handler := &getMovieDetailHandler{
				interactor: movie.NewGetMovieDetailInteractor(
					svc.MovieService(),
					svc.ReviewService(),
					tmdbSvc,
				),
			}

			// Act
			err := handler.handle(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusOK, rec.Code)
			assert.JSONEq(t, tt.expected, rec.Body.String())
		})
	}
}
