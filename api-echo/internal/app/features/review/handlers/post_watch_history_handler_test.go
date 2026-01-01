package handlers

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/review"
	"filmogophery/internal/app/responses"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/config"
	"filmogophery/internal/pkg/constant"
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/tests"
)

// 異常系
func TestPostReviewHistoryHandler_handle__Error(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	var operator *model.Users
	tx.First(&operator)

	// Create movies
	m := tests.CreateMovies(t, tx, &model.Movies{
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 28, 11, 22, 33, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})
	// Create reviews
	rv := tests.CreateReviews(t, tx, &model.Reviews{
		ID: 314, UserID: 1, MovieID: m.ID,
	})
	tomorrowDate := time.Now().AddDate(0, 0, 1).Format(constant.DateFormat)

	for _, tt := range []struct {
		testCase       string
		reviewID       int32
		body           string
		expectedCode   int
		expectedErrors map[string][]string
	}{
		{
			testCase:       "存在しないレビューIDを指定",
			reviewID:       404,
			body:           `{"platformId": 99, "watchedDate": "2025-10-28"}`,
			expectedCode:   http.StatusNotFound,
			expectedErrors: map[string][]string{"id": {"404"}},
		},
		{
			testCase:       "存在しないプラットフォームIDを指定",
			reviewID:       rv.ID,
			body:           `{"platformId": 404, "watchedDate": "2025-10-28"}`,
			expectedCode:   http.StatusNotFound,
			expectedErrors: map[string][]string{"id": {"404"}},
		},
		{
			testCase:       "不正な視聴日付を指定_スラッシュ区切り",
			reviewID:       rv.ID,
			body:           `{"platformId": 99, "watchedDate": "2025/10/28"}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"failed to parse date(2025/10/28)"}},
		},
		{
			testCase:       "不正な視聴日付を指定_時刻付き",
			reviewID:       rv.ID,
			body:           `{"platformId": 99, "watchedDate": "2025-10-28T10:20:30"}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"failed to parse date(2025-10-28T10:20:30)"}},
		},
		{
			testCase:       "不正な視聴日付を指定_存在しない月",
			reviewID:       rv.ID,
			body:           `{"platformId": 99, "watchedDate": "2025-13-28"}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"failed to parse date(2025-13-28)"}},
		},
		{
			testCase:       "不正な視聴日付を指定_存在しない日",
			reviewID:       rv.ID,
			body:           `{"platformId": 99, "watchedDate": "2025-02-30"}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"failed to parse date(2025-02-30)"}},
		},
		{
			testCase:       "不正な視聴日付を指定_0埋めなし",
			reviewID:       rv.ID,
			body:           `{"platformId": 99, "watchedDate": "2025-1-1"}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"failed to parse date(2025-1-1)"}},
		},
		{
			testCase:       "不正な視聴日付を指定_空文字",
			reviewID:       rv.ID,
			body:           `{"platformId": 99, "watchedDate": ""}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"failed to parse date()"}},
		},
		{
			testCase:       "未来の視聴日付を指定",
			reviewID:       rv.ID,
			body:           fmt.Sprintf(`{"platformId": 99, "watchedDate": "%s"}`, tomorrowDate),
			expectedCode:   http.StatusBadRequest,
			expectedErrors: map[string][]string{"WatchedDate": {"date cannot be in the future"}},
		},
	} {
		t.Run(tt.testCase, func(t *testing.T) {
			reviewID := strconv.Itoa(int(tt.reviewID))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/reviews/"+reviewID+"/history", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(reviewID)

			c.Set("operator", operator)

			svc := services.NewServiceContainer(tx, conf, nil, nil, nil)
			handler := &postReviewHistoryHandler{
				interactor: review.NewCreateReviewHistoryInteractor(
					svc.ReviewService(),
					svc.PlatformService(),
				),
			}

			// Act
			err := handler.handle(c)

			// Assert
			assert.Error(t, err)
			he, ok := err.(*echo.HTTPError)
			assert.True(t, ok, "error should be *echo.HTTPError")
			assert.Equal(t, tt.expectedCode, he.Code)

			errResp := he.Message.(responses.ErrorResponse)
			assert.NotEmpty(t, errResp.Message)
			assert.Equal(t, tt.expectedErrors, errResp.Errors)
		})
	}
}

// 正常系
func TestPostReviewHistoryHandler_handle(t *testing.T) {
	for _, tt := range []struct {
		testCase          string
		body              string
		expectWatchedDate time.Time
	}{
		{
			testCase:          "視聴日付の指定なし",
			body:              `{"platformId": 99}`,
			expectWatchedDate: time.Date(1895, 12, 28, 0, 0, 0, 0, time.Local), // デフォルトの日付が設定されること
		},
		{
			testCase:          "視聴日付の指定あり",
			body:              `{"platformId": 99, "watchedDate": "2025-10-28"}`,
			expectWatchedDate: time.Date(2025, 10, 28, 0, 0, 0, 0, time.Local),
		},
	} {
		t.Run(tt.testCase, func(t *testing.T) {
			// Arrange
			testDB := tests.SetupTestDB()
			conf := config.LoadConfig()
			tx := testDB.Begin()
			defer tx.Rollback()

			var operator *model.Users
			tx.First(&operator)

			// Create movies
			m := tests.CreateMovies(t, tx, &model.Movies{
				ID:             1,
				Title:          "テスト映画タイトル1",
				Overview:       "テスト概要1",
				ReleaseDate:    time.Date(2025, 10, 28, 11, 22, 33, 456789, time.Local),
				RuntimeMinutes: 314,
				TmdbID:         1592,
			})
			// Create reviews
			rv := tests.CreateReviews(t, tx, &model.Reviews{
				ID: 314, UserID: 1, MovieID: m.ID,
			})

			// テスト前の視聴履歴数を確認
			var beforeCount int64
			tx.Model(&model.WatchHistory{}).Where("movie_id = ?", m.ID).Count(&beforeCount)

			reviewID := strconv.Itoa(int(rv.ID))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/reviews/"+reviewID+"/history", strings.NewReader(tt.body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(reviewID)

			c.Set("operator", operator)

			svc := services.NewServiceContainer(tx, conf, nil, nil, nil)
			handler := &postReviewHistoryHandler{
				interactor: review.NewCreateReviewHistoryInteractor(
					svc.ReviewService(),
					svc.PlatformService(),
				),
			}

			// Act
			err := handler.handle(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, rec.Code)

			// データベースの状態確認
			var afterCount int64
			tx.Model(&model.WatchHistory{}).Where("movie_id = ?", m.ID).Count(&afterCount)
			assert.Equal(t, beforeCount+1, afterCount)

			// 作成された視聴履歴の内容確認
			var createdWatchHistory model.WatchHistory
			err = tx.Where("movie_id = ?", m.ID).First(&createdWatchHistory).Error
			assert.NoError(t, err)
			assert.Equal(t, int32(1), createdWatchHistory.UserID)
			assert.Equal(t, m.ID, createdWatchHistory.MovieID)
			assert.Equal(t, int32(99), createdWatchHistory.PlatformID)
			if createdWatchHistory.WatchedDate != nil {
				assert.True(t, tt.expectWatchedDate.Equal(*createdWatchHistory.WatchedDate))
			} else {
				assert.Nil(t, createdWatchHistory.WatchedDate)
			}
		})
	}
}
