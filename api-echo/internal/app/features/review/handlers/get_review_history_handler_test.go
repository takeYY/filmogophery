// TODO: /movies/{movieId}/watch-history に変更するため、いずれ消す

package handlers

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/review"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/tests"
)

// 異常系
func TestGetWatchHistoryHandler_handle__Error(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	var operator *model.Users
	tx.First(&operator)
	id := "404"

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/reviews/"+id+"/history", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues(id)

	c.Set("operator", operator)

	svc := services.NewServiceContainer(tx, conf, nil, nil, nil)
	handler := &getReviewHistoryHandler{
		interactor: review.NewGetReviewHistoryInteractor(
			svc.ReviewService(),
		),
	}

	// Act
	err := handler.handle(c)

	// Assert
	assert.Error(t, err)
	he, ok := err.(*echo.HTTPError)
	assert.True(t, ok, "error should be *echo.HTTPError")
	assert.Equal(t, http.StatusNotFound, he.Code)

	errResp := he.Message.(responses.ErrorResponse)
	assert.Equal(t, "review not found", errResp.Message)
	assert.Equal(t, map[string][]string{"id": {"404"}}, errResp.Errors)
}

// 正常系
func TestGetWatchHistoryHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	var operator *model.Users
	tx.First(&operator)

	// Create movies
	m1 := tests.CreateMovies(t, tx, &model.Movies{
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 20, 11, 22, 33, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})
	m2 := tests.CreateMovies(t, tx, &model.Movies{
		ID:             2,
		Title:          "テスト映画タイトル2",
		Overview:       "テスト概要2",
		ReleaseDate:    time.Date(2025, 10, 20, 11, 22, 33, 456789, time.Local),
		RuntimeMinutes: 65,
		TmdbID:         35,
	})
	// Create reviews
	rv1 := tests.CreateReviews(t, tx, &model.Reviews{
		ID: 314, UserID: 1, MovieID: m1.ID,
	})
	rv2 := tests.CreateReviews(t, tx, &model.Reviews{
		ID: 1592, UserID: 1, MovieID: m2.ID,
	})
	// Create watch_history
	tests.CreateWatchHistory(t, tx, &model.WatchHistory{
		ID: 65, UserID: 1, MovieID: 2, PlatformID: 99, WatchedDate: &[]time.Time{time.Date(2025, 1, 2, 0, 0, 0, 0, time.Local)}[0],
	})
	tests.CreateWatchHistory(t, tx, &model.WatchHistory{
		ID: 8979, UserID: 1, MovieID: 2, PlatformID: 1, WatchedDate: &[]time.Time{time.Date(2025, 2, 3, 0, 0, 0, 0, time.Local)}[0],
	})

	for _, tt := range []struct {
		testCase string
		reviewID int32
		expected string
	}{
		{
			testCase: "視聴履歴がないレビュー",
			reviewID: rv1.ID,
			expected: `[]`,
		},
		{
			testCase: "視聴履歴が複数あるレビュー",
			reviewID: rv2.ID,
			expected: `[
				{
					"id": 65,
					"platform": {
						"code": "unknown",
						"name": "不明"
					},
					"watchedAt": "2025-01-02T00:00:00+09:00"
				},
				{
					"id": 8979,
					"platform": {
						"code": "primeVideo",
						"name": "Prime Video"
					},
					"watchedAt": "2025-02-03T00:00:00+09:00"
				}
			]`,
		},
	} {
		t.Run(tt.testCase, func(t *testing.T) {
			reviewID := strconv.Itoa(int(tt.reviewID))

			e := echo.New()
			req := httptest.NewRequest(http.MethodGet, "/reviews/"+reviewID+"/history", nil)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(reviewID)

			c.Set("operator", operator)

			svc := services.NewServiceContainer(tx, conf, nil, nil, nil)
			handler := &getReviewHistoryHandler{
				interactor: review.NewGetReviewHistoryInteractor(
					svc.ReviewService(),
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
