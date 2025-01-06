package health

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/config"
	"filmogophery/internal/health"
)

func TestHealth(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	c := e.NewContext(req, rec)
	c.SetPath("/health")

	// 設定ファイルの読み込み
	conf := config.LoadConfig()
	handler := health.NewHandler(conf)
	er := handler.ReaderHandler.Health(c)
	if er != nil {
		t.Fatal(er)
	}

	if assert.NoError(t, er) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expected := `{"message":"system all green"}`
		assert.JSONEq(t, expected, rec.Body.String())
	}
}
