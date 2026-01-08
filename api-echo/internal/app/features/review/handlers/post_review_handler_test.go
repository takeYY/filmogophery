package handlers

import (
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
	"filmogophery/internal/pkg/gen/model"
	"filmogophery/tests"
)

// 異常系
func TestPostReviewHandler_handle__Error(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	tx := testDB.Begin()
	defer tx.Rollback()

	var operator *model.Users
	tx.First(&operator)

	m := tests.CreateMovies(t, tx, &model.Movies{
		ID:             1,
		Title:          "テスト映画タイトル1",
		Overview:       "テスト概要1",
		ReleaseDate:    time.Date(2025, 10, 15, 10, 20, 30, 456789, time.Local),
		RuntimeMinutes: 314,
		TmdbID:         1592,
	})

	for _, tt := range []struct {
		testCase       string
		movieID        int32
		reqBody        string
		expectedStatus int
		expectedErrors map[string][]string
	}{
		{
			testCase:       "ratingとcommentの両方がnull",
			movieID:        m.ID,
			reqBody:        `{"rating": null, "comment": null}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: map[string][]string{
				"Rating":  {"both rating and comment cannot be null"},
				"Comment": {"both rating and comment cannot be null"},
			},
		},
		{
			testCase:       "ratingが最小値未満",
			movieID:        m.ID,
			reqBody:        `{"rating": 0.0}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: map[string][]string{
				"Rating": {"Rating validation failed on gte"},
			},
		},
		{
			testCase:       "ratingが最大値超過",
			movieID:        m.ID,
			reqBody:        `{"rating": 5.1}`,
			expectedStatus: http.StatusBadRequest,
			expectedErrors: map[string][]string{
				"Rating": {"Rating validation failed on lte"},
			},
		},
		{
			testCase:       "存在しない映画ID",
			movieID:        404,
			reqBody:        `{"rating": 3.1, "comment": "something stylish comments"}`,
			expectedStatus: http.StatusNotFound,
			expectedErrors: map[string][]string{
				"id": {"404"},
			},
		},
	} {
		t.Run(tt.testCase, func(t *testing.T) {
			id := strconv.Itoa(int(tt.movieID))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/movies/"+id+"/reviews", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(id)

			c.Set("operator", operator)

			svc := services.NewServiceContainer(tx, conf, nil, nil, nil)
			handler := &postReviewHandler{
				interactor: review.NewCreateReviewInteractor(
					svc.MovieService(),
					svc.ReviewService(),
				),
			}

			// Act
			err := handler.handle(c)

			// Assert
			assert.Error(t, err)
			he, ok := err.(*echo.HTTPError)
			assert.True(t, ok, "error should be *echo.HTTPError")
			assert.Equal(t, tt.expectedStatus, he.Code)

			errResp := he.Message.(responses.ErrorResponse)
			assert.NotEmpty(t, errResp.Message)
			assert.Equal(t, tt.expectedErrors, errResp.Errors)
		})
	}
}

// 正常系
func TestPostReviewHandler_handle(t *testing.T) {
	for _, tt := range []struct {
		testCase        string
		reqBody         string
		expectedRating  *float64
		expectedComment *string
	}{
		{
			testCase:        "ratingとcommentが両方ある",
			reqBody:         `{"rating": 3.1, "comment": "something stylish comments"}`,
			expectedRating:  &[]float64{3.1}[0],
			expectedComment: &[]string{"something stylish comments"}[0],
		},
		{
			testCase:        "ratingのみ",
			reqBody:         `{"rating": 4.5}`,
			expectedRating:  &[]float64{4.5}[0],
			expectedComment: nil,
		},
		{
			testCase:        "commentのみ",
			reqBody:         `{"comment": "great movie"}`,
			expectedRating:  nil,
			expectedComment: &[]string{"great movie"}[0],
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

			m := tests.CreateMovies(t, tx, &model.Movies{
				ID:             65,
				Title:          "テスト映画タイトル1",
				Overview:       "テスト概要1",
				ReleaseDate:    time.Date(2025, 10, 15, 10, 20, 30, 456789, time.Local),
				RuntimeMinutes: 314,
				TmdbID:         1592,
			})

			// テスト前のレビュー数を確認
			var beforeCount int64
			tx.Model(&model.Reviews{}).Where("movie_id = ?", m.ID).Count(&beforeCount)

			id := strconv.Itoa(int(m.ID))

			e := echo.New()
			req := httptest.NewRequest(http.MethodPost, "/movies/"+id+"/reviews", strings.NewReader(tt.reqBody))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)
			c.SetParamNames("id")
			c.SetParamValues(id)

			c.Set("operator", operator)

			svc := services.NewServiceContainer(tx, conf, nil, nil, nil)
			handler := &postReviewHandler{
				interactor: review.NewCreateReviewInteractor(
					svc.MovieService(),
					svc.ReviewService(),
				),
			}

			// Act
			err := handler.handle(c)

			// Assert
			assert.NoError(t, err)
			assert.Equal(t, http.StatusCreated, rec.Code)
			assert.Empty(t, rec.Body.String())

			// データベースの状態確認
			var afterCount int64
			tx.Model(&model.Reviews{}).Where("movie_id = ?", m.ID).Count(&afterCount)
			assert.Equal(t, beforeCount+1, afterCount)

			// 作成されたレビューの内容確認
			var createdReview model.Reviews
			err = tx.Where("movie_id = ? AND user_id = ?", m.ID, 1).First(&createdReview).Error
			assert.NoError(t, err)
			assert.Equal(t, int32(1), createdReview.UserID)
			assert.Equal(t, m.ID, createdReview.MovieID)
			assert.Equal(t, tt.expectedRating, createdReview.Rating)
			assert.Equal(t, tt.expectedComment, createdReview.Comment)
			assert.NotZero(t, createdReview.CreatedAt)
			assert.NotZero(t, createdReview.UpdatedAt)
		})
	}
}
