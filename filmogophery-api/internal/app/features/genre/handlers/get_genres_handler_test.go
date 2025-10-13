package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/genre"
	"filmogophery/internal/app/services"
	"filmogophery/internal/config"
	"filmogophery/tests"
)

func TestGetGenresHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	svc := services.NewServiceContainer(testDB, conf)

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
		{"code": "action", "name": "アクション"},
		{"code": "adventure", "name": "アドベンチャー"},
		{"code": "animation", "name": "アニメーション"},
		{"code": "comedy", "name": "コメディ"},
		{"code": "crime", "name": "クライム"},
		{"code": "documentary", "name": "ドキュメンタリー"},
		{"code": "drama", "name": "ドラマ"},
		{"code": "family", "name": "ファミリー"},
		{"code": "fantasy", "name": "ファンタジー"},
		{"code": "history", "name": "ヒストリー"},
		{"code": "horror", "name": "ホラー"},
		{"code": "musical", "name": "ミュージカル"},
		{"code": "mystery", "name": "ミステリー"},
		{"code": "romance", "name": "ロマンス"},
		{"code": "sf", "name": "SF"},
		{"code": "tv", "name": "TV"},
		{"code": "thriller", "name": "スリラー"},
		{"code": "war", "name": "戦争"},
		{"code": "western", "name": "西部劇"}
	]`

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
