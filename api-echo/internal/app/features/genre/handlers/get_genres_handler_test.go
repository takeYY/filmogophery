package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/genre"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/config"
	"filmogophery/tests"
)

func TestGetGenresHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	svc := services.NewServiceContainer(testDB, conf, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &getGenresHandler{
		interactor: genre.NewGetGenresInteractor(svc.GenreService()),
	}

	// Act
	err := handler.handle(c)

	// expected
	expected := `[
		{"code": "adventure", "name": "アドベンチャー"},
		{"code": "fantasy", "name": "ファンタジー"},
		{"code": "animation", "name": "アニメーション"},
		{"code": "drama", "name": "ドラマ"},
		{"code": "horror", "name": "ホラー"},
		{"code": "action", "name": "アクション"},
		{"code": "comedy", "name": "コメディ"},
		{"code": "history", "name": "ヒストリー"},
		{"code": "western", "name": "西部劇"},
		{"code": "thriller", "name": "スリラー"},
		{"code": "crime", "name": "クライム"},
		{"code": "documentary", "name": "ドキュメンタリー"},
		{"code": "sf", "name": "SF"},
		{"code": "mystery", "name": "ミステリー"},
		{"code": "musical", "name": "ミュージカル"},
		{"code": "romance", "name": "ロマンス"},
		{"code": "family", "name": "ファミリー"},
		{"code": "war", "name": "戦争"},
		{"code": "tv", "name": "TV"}
	]`

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
