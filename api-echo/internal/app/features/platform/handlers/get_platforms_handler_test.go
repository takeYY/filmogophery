package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"

	"filmogophery/internal/app/features/platform"
	"filmogophery/internal/app/services"
	"filmogophery/internal/pkg/config"
	"filmogophery/tests"
)

func TestGetPlatformsHandler_handle(t *testing.T) {
	// Arrange
	testDB := tests.SetupTestDB()
	conf := config.LoadConfig()
	svc := services.NewServiceContainer(testDB, conf, nil, nil, nil)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/genres", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	handler := &getPlatformsHandler{
		interactor: platform.NewGetPlatformsInteractor(svc.PlatformService()),
	}

	// Act
	err := handler.handle(c)

	// expected
	expected := `[
		{"code": "primeVideo", "name": "Prime Video"},
		{"code": "netflix", "name": "Netflix"},
		{"code": "uNext", "name": "U-NEXT"},
		{"code": "disneyPlus", "name": "Disney+"},
		{"code": "youtube", "name": "YouTube"},
		{"code": "appleTv", "name": "Apple TV+"},
		{"code": "hulu", "name": "Hulu"},
		{"code": "dAnime", "name": "dアニメ"},
		{"code": "telasa", "name": "TELASA"},
		{"code": "cinema", "name": "映画館"},
		{"code": "unknown", "name": "不明"}
	]`

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.JSONEq(t, expected, rec.Body.String())
}
